<!-- src/App.svelte -->
<script lang="ts">
import { api } from '$lib/api';
import DeviceList from '$lib/components/DeviceList.svelte';
import DeviceModal from '$lib/components/DeviceModal.svelte';
import SettingsModal from '$lib/components/SettingsModal.svelte';
import Toast from '$lib/components/Toast.svelte';
import { devices } from '$lib/store/devices.svelte';
import type { Device } from '$lib/types';
import { onMount } from 'svelte';

// State using Svelte 5 runes
let currentlyConnectedId = $state<string | null>(null);
let isLoading = $derived(devices.isLoading);
let error = $state<string | null>(null);
let showDeviceModal = $state(false);
let showSettingsModal = $state(false);
let editingDevice = $state<Device | null>(null);
let toastMessage = $state<string | null>(null);

// Load data on component mount
onMount(async () => {
  // Check for current connection status
  // In a real implementation, we would need a way to get the current connection status from the backend
  // For now, we'll simulate this with a simple message handler
  window.addEventListener('message', event => {
    if (event.data.type === 'connection-status') {
      currentlyConnectedId = event.data.deviceId;
    }
  });

  window.addEventListener('keyup', e => {
    if (e.key === 'Esc' || e.key === 'Escape') {
      showSettingsModal = false;
      showDeviceModal = false;
    }
  });
});

// Handle device actions
async function handleConnect(device: Device) {
  try {
    if (currentlyConnectedId === device.id) {
      await api.disconnect();
      toastMessage = `Disconnected from ${device.name}`;
      currentlyConnectedId = null;
    } else {
      await api.connectToDevice(device.id);
      toastMessage = `Connected to ${device.name}`;
      currentlyConnectedId = device.id;
    }
  } catch (err) {
    error = err instanceof Error ? err.message : 'Connection failed';
  }
}

async function handleShowEditModal(device: Device) {
  editingDevice = device;
  showDeviceModal = true;
}

async function handleDelete(deviceId: string) {
  if (!confirm('Are you sure you want to delete this device?')) {
    return;
  }

  try {
    // devices = devices.filter(d => d.id !== deviceId);
    await devices.delete(deviceId);
    toastMessage = 'Device deleted successfully';

    if (currentlyConnectedId === deviceId) {
      currentlyConnectedId = null;
    }
  } catch (err) {
    error = err instanceof Error ? err.message : 'Failed to delete device';
  }
}

function handleShowAddModal() {
  editingDevice = null;
  showDeviceModal = true;
}

function handleShowSettingsModal() {
  showSettingsModal = true;
}

async function handleDeviceSave(device: Device) {
  try {
    if (device.id) {
      // Update existing device
      await devices.update(device);
      toastMessage = 'Device updated successfully';
    } else {
      // Add a new device
      await devices.add(device);
      toastMessage = 'Device added successfully';
    }
    showDeviceModal = false;
  } catch (err) {
    error = err instanceof Error ? err.message : 'Failed to save device';
  }
}

function handleModalClose() {
  showDeviceModal = false;
  showSettingsModal = false;
}

function clearError() {
  error = null;
}
</script>

<main class="container">
    <header>
        <h1>Heimdall</h1>
        <h2 class="subtitle">Remote View Connection Manager</h2>
        <div class="actions">
            <button onclick={handleShowAddModal} class="btn btn-success">Add Device</button>
            <button onclick={handleShowSettingsModal} class="btn btn-secondary">Settings</button>
        </div>
    </header>

    <!-- Status Bar -->
    <!--    <div class="status-bar {currentlyConnectedId ? 'connected' : ''}">-->
    <div class={['status-bar', {connected: currentlyConnectedId }]}>
        {#if currentlyConnectedId}
      <span>
        Connected to:
          {devices.search(currentlyConnectedId)?.name || 'Unknown Device'}
      </span>
            <button onclick={() => api.disconnect()} class="btn btn-danger">Disconnect</button>
        {:else}
            <span>Not connected to any device</span>
        {/if}
    </div>

    <!-- Loading indicator -->
    {#if isLoading}
        <div class="loading">Loading devices...</div>
    {/if}

    <!-- Error message -->
    {#if error}
        <div class="error">
            {error}
            <button onclick={clearError} class="btn-close">&times;</button>
        </div>
    {/if}

    <!-- Device list -->
    {#if devices.filtered.length > 0}
        <DeviceList
                connectedId={currentlyConnectedId}
                onConnect={handleConnect}
                onEdit={handleShowEditModal}
                onDelete={handleDelete}
        />
    {:else if !isLoading}
        <div class="empty-state">
            <p>No devices configured yet. Add your first device to get started.</p>
            <button onclick={handleShowAddModal} class="btn btn-primary">Add Device</button>
        </div>
    {/if}
    <!-- Modals -->
    <DeviceModal
            bind:showModal={showDeviceModal}
            device={editingDevice}
            onSave={handleDeviceSave}
            onClose={handleModalClose}
    />

    <SettingsModal
            bind:showModal={showSettingsModal}
            onClose={handleModalClose}
    />

    <!-- Toast notifications -->
    {#if toastMessage}
        <Toast message={toastMessage}/>
    {/if}
</main>

<style>
    .container {
        max-width: var(--heimdall-container-width);
        margin: 0 auto;
    }

    header {
        display: block;
        justify-content: space-between;
        align-items: center;
        margin-bottom: var(--heimdall-spacing-xl);

        h1, h2 {
            margin: 0;
            padding: 0;
            line-height: 1;
        }

        h2.subtitle {
            margin-bottom: var(--heimdall-spacing-lg);
        }

        .actions {
            margin: var(--heimdall-spacing-sm) 0 0 0;
            display: flex;
            gap: var(--heimdall-spacing-sm);
            flex-wrap: wrap;
        }
    }

    .status-bar {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: var(--heimdall-spacing-md) var(--heimdall-spacing-lg);
        background-color: var(--heimdall-bg-card);
        border: 1px solid var(--heimdall-border-color);
        border-radius: var(--heimdall-rounded);
        margin-bottom: var(--heimdall-spacing-lg);

        &.connected {
            background-color: var(--heimdall-connected-bg);
            border-color: var(--heimdall-connected-border);
        }
    }

    .loading {
        text-align: center;
        padding: var(--heimdall-spacing-2xl);
        color: var(--heimdall-text-heading);
    }

    .error {
        background-color: var(--heimdall-error-bg);
        border: 1px solid var(--heimdall-error-border);
        color: var(--heimdall-error-text);
        padding: var(--heimdall-spacing-md) var(--heimdall-spacing-lg);
        border-radius: var(--heimdall-rounded);
        margin-bottom: var(--heimdall-spacing-xl);
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .btn-close {
        background: none;
        border: none;
        color: var(--heimdall-error-text);
        font-size: var(--heimdall-font-size-xl);
        cursor: pointer;
    }

    .empty-state {
        text-align: center;
        padding: var(--heimdall-spacing-3xl);
        background-color: var(--heimdall-bg-card);
        border: 1px solid var(--heimdall-border-color);
        border-radius: var(--heimdall-rounded);

        p {
            margin: 0 0 var(--heimdall-spacing-lg) 0;
        }
    }
</style>
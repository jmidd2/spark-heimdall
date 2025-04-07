<!-- src/App.svelte -->
<script lang="ts">
  import { api } from '$lib/api';
  import { onMount } from 'svelte';
  import type { Device } from '$lib/types';
  import DeviceList from '$lib/components/DeviceList.svelte';
  import DeviceModal from '$lib/components/DeviceModal.svelte';
  import SettingsModal from '$lib/components/SettingsModal.svelte';
  import Toast from '$lib/components/Toast.svelte';

  // State using Svelte 5 runes
  let devices = $state<Device[]>([]);
  let currentlyConnectedId = $state<string | null>(null);
  let isLoading = $state(true);
  let error = $state<string | null>(null);
  let showDeviceModal = $state(false);
  let showSettingsModal = $state(false);
  let editingDevice = $state<Device | null>(null);
  let toastMessage = $state<string | null>(null);

  $inspect(devices)

  // Load data on component mount
  onMount(async () => {
    try {
      isLoading = true;
      devices = await api.getDevices();

      // Check for current connection status
      // In a real implementation, we would need a way to get the current connection status from the backend
      // For now, we'll simulate this with a simple message handler
      window.addEventListener('message', (event) => {
        if (event.data.type === 'connection-status') {
          currentlyConnectedId = event.data.deviceId;
        }
      });
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to load devices';
    } finally {
      isLoading = false;
    }
  });

  // Show toast message for 3 seconds
  $effect(() => {
    if (toastMessage) {
      const timeout = setTimeout(() => {
        toastMessage = null;
      }, 3000);

      return () => clearTimeout(timeout);
    }
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

  async function handleEdit(device: Device) {
    editingDevice = device;
    showDeviceModal = true;
  }

  async function handleDelete(deviceId: string) {
    if (!confirm('Are you sure you want to delete this device?')) {
      return;
    }

    try {
      await api.deleteDevice(deviceId);
      devices = devices.filter(d => d.id !== deviceId);
      toastMessage = 'Device deleted successfully';

      if (currentlyConnectedId === deviceId) {
        currentlyConnectedId = null;
      }
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to delete device';
    }
  }

  function handleAddDevice() {
    editingDevice = null;
    showDeviceModal = true;
  }

  function handleOpenSettings() {
    showSettingsModal = true;
  }

  async function handleDeviceSave(device: Device) {
    try {
      if (device.id) {
        // Update existing device
        const updated = await api.updateDevice(device);
        devices = devices.map(d => d.id === device.id ? updated : d);
        toastMessage = 'Device updated successfully';
      } else {
        // Add new device
        const added = await api.addDevice(device);
        devices = [...devices, added];
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
    <h1>Heimdall - Remote Connection Manager</h1>
    <div class="actions">
      <button onclick={handleAddDevice} class="btn btn-primary">Add Device</button>
      <button onclick={handleOpenSettings} class="btn btn-secondary">Settings</button>
    </div>
  </header>

  <!-- Status Bar -->
  <div class="status-bar {currentlyConnectedId ? 'connected' : ''}">
    {#if currentlyConnectedId}
      <span>
        Connected to:
        {devices.find(d => d.id === currentlyConnectedId)?.name || 'Unknown Device'}
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
  {#if devices.length > 0}
    <DeviceList
            {devices}
            connectedId={currentlyConnectedId}
            onConnect={handleConnect}
            onEdit={handleEdit}
            onDelete={handleDelete}
    />
  {:else if !isLoading}
    <div class="empty-state">
      <p>No devices configured yet. Add your first device to get started.</p>
      <button onclick={handleAddDevice} class="btn btn-primary">Add Device</button>
    </div>
  {/if}

  <!-- Modals -->
  {#if showDeviceModal}
    <DeviceModal
            device={editingDevice}
            onSave={handleDeviceSave}
            onClose={handleModalClose}
    />
  {/if}

  {#if showSettingsModal}
    <SettingsModal
            devices={devices}
            onClose={handleModalClose}
    />
  {/if}

  <!-- Toast notifications -->
  {#if toastMessage}
    <Toast message={toastMessage} />
  {/if}
</main>

<style>
  .container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 1rem;
  }

  header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }

  .actions {
    display: flex;
    gap: 0.5rem;
  }

  .status-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem 1rem;
    background-color: #1b1e1f;
    border: 1px solid #262a2b;
    border-radius: 4px;
    margin-bottom: 1.5rem;
  }

  .status-bar.connected {
    background-color: #1a3e29;
    border-color: #245931;
  }

  .loading {
    text-align: center;
    padding: 2rem;
    color: #aec2d3;
  }

  .error {
    background-color: #482121;
    border: 1px solid #692929;
    color: #f5c2c2;
    padding: 0.75rem 1rem;
    border-radius: 4px;
    margin-bottom: 1.5rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .btn-close {
    background: none;
    border: none;
    color: #f5c2c2;
    font-size: 1.25rem;
    cursor: pointer;
  }

  .empty-state {
    text-align: center;
    padding: 3rem;
    background-color: #1b1e1f;
    border: 1px solid #262a2b;
    border-radius: 4px;
  }
</style>
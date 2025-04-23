<!-- frontend/src/lib/components/SettingsModal.svelte -->
<script lang="ts">
import { api } from '$lib/api';
// biome-ignore lint/correctness/noUnusedImports: <explanation>
import Modal from '$lib/components/Modal.svelte';
import type { AppConfig, Device } from '$lib/types';
import { onMount } from 'svelte';

// Props
type SettingsModalProps = {
  devices: Device[];
  onClose: () => void;
  showModal: boolean;
};
let {
  devices,
  onClose,
  showModal = $bindable(),
}: SettingsModalProps = $props();

// Local state
let formData = $state<AppConfig>({
  server: {
    port: 8080,
  },
  connection: {
    auto_start: false,
    auto_start_id: '',
  },
  clients: {
    vnc_viewer: '',
    vnc_password_file: '',
    rdp_viewer: '',
  },
  logging: {
    level: 'info',
    format: 'text',
  },
});

let isLoading = $state(true);
let error = $state<string | null>(null);
let success = $state<boolean | null>(null);

// Load current settings when component mounts
onMount(async () => {
  try {
    isLoading = true;
    const config = await api.getConfig();
    formData.server.port = config.server.port;
    formData.connection.auto_start = config.connection.auto_start;
    formData.connection.auto_start_id = config.connection.auto_start_id;
    formData.clients.vnc_viewer = config.clients.vnc_viewer;
    formData.clients.vnc_password_file = config.clients.vnc_password_file;
    formData.clients.rdp_viewer = config.clients.rdp_viewer;
  } catch (err) {
    error = err instanceof Error ? err.message : 'Failed to load settings';
  } finally {
    isLoading = false;
  }
});

$effect(() => {
  if (success) {
    const timeout = setTimeout(() => {
      success = false;
    }, 3000);

    return () => clearTimeout(timeout);
  }
});

// Form validation
function validateForm(): boolean {
  if (formData.server.port <= 0 || formData.server.port > 65535) {
    error = 'Port must be between 1 and 65535';
    return false;
  }

  error = null;
  return true;
}

// Save settings
async function handleSubmit() {
  if (!validateForm()) {
    return;
  }

  try {
    const result = await api.updateConfig(formData);
    if (result.success) {
      success = true;
    } else {
      error = 'Failed to save settings';
    }
  } catch (err) {
    console.error(err);
    error = err instanceof Error ? err.message : 'Failed to save settings';
  }
}

function handleRestartServer() {
  // This would typically send a request to restart the server
  // For now, we'll just show a message that this feature is not implemented
  alert('Server restart functionality is not implemented yet.');
}
</script>

<Modal title="Application Settings" showModal={showModal} onSubmit={handleSubmit} onClose={onClose}>
    {#if isLoading}
        <div class="loading">Loading settings...</div>
    {:else}
        <div class="form-section">
            <div class="form-row">
                <h3>Server Settings</h3>

                <div class="form-group">
                    <label for="port">HTTP Port</label>
                    <input
                            type="number"
                            id="port"
                            bind:value={formData.server.port}
                            min="1"
                            max="65535"
                    />
                    <small>Changes require server restart</small>
                </div>
            </div>

            <div class="form-row">
                <h3>External Applications</h3>

                <div class="form-group">
                    <label for="vnc_viewer">VNC Viewer Path</label>
                    <input
                            type="text"
                            id="vnc_viewer"
                            bind:value={formData.clients.vnc_viewer}
                            placeholder="vncviewer"
                    />
                </div>

                <div class="form-group">
                    <label for="vnc_password_file">VNC Password File</label>
                    <input
                            type="text"
                            id="vnc_password_file"
                            bind:value={formData.clients.vnc_password_file}
                            placeholder="~/.vnc/passwd"
                    />
                </div>

                <div class="form-group">
                    <label for="rdp_viewer">RDP Viewer Path</label>
                    <input
                            type="text"
                            id="rdp_viewer"
                            bind:value={formData.clients.rdp_viewer}
                            placeholder="xfreerdp"
                    />
                </div>
            </div>

            <div class="form-row">
                <h3>Connection Settings</h3>

                <div class="form-group checkbox-group">
                    <input
                            type="checkbox"
                            id="auto_start"
                            bind:checked={formData.connection.auto_start}
                    />
                    <div>
                        <label for="auto_start">Auto-start connection on launch</label>
                    </div>
                </div>

                <div class="form-group">
                    <label for="auto_start_id">Auto-start Device</label>
                    <select
                            id="auto_start_id"
                            bind:value={formData.connection.auto_start_id}
                            disabled={!formData.connection.auto_start}
                    >
                        <option value="">None</option>
                        {#each devices as device}
                            <option value={device.id}>{device.name}</option>
                        {/each}
                    </select>
                </div>
            </div>
            <div class="form-row">
                <h3>Actions</h3>

                <button type="button" class="btn btn-danger" onclick={handleRestartServer}>
                    Restart Server
                </button>
            </div>
        </div>

        {#if error}
            <div class="error-message">{error}</div>
        {/if}
        {#if success}
            <div class="success-message">Settings saved!</div>
        {/if}
    {/if}
</Modal>

<style>

</style>
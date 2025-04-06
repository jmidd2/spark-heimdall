<!-- frontend/src/lib/components/SettingsModal.svelte -->
<script lang="ts">
  import { api } from '$lib/api';
  import type { Device, AppConfig } from '$lib/types';
  import { onMount } from 'svelte';

  // Props
  type SettingsModalProps = {
    devices: Device[];
    onClose: () => void;
  }
  const { devices, onClose }: SettingsModalProps = $props()

  // Local state
  let formData = $state<AppConfig>({
    server: {
      port: 8080
    },
    connection: {
      auto_start: false,
      auto_start_id: ''
    },
    clients: {
      vnc_viewer: '',
      vnc_password_file: '',
      rdp_viewer: ''
    },
    logging: {
      level: 'info',
      format: 'text'
    }
  });

  let isLoading = $state(true);
  let error = $state<string | null>(null);
  let isSaving = $state(false);
  let saveSuccess = $state(false);

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

  // Show success message for 3 seconds
  $effect(() => {
    if (saveSuccess) {
      const timeout = setTimeout(() => {
        saveSuccess = false;
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
  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    try {
      isSaving = true;
      await api.updateConfig(formData);
      saveSuccess = true;
    } catch (err) {
      error = err instanceof Error ? err.message : 'Failed to save settings';
    } finally {
      isSaving = false;
    }
  }

  function handleRestartServer() {
    // This would typically send a request to restart the server
    // For now, we'll just show a message that this feature is not implemented
    alert('Server restart functionality is not implemented yet.');
  }
</script>

<div class="modal-backdrop" onclick={onClose}>
    <div class="modal-content">
        <div class="modal-header">
            <h2>Application Settings</h2>
            <button type="button" class="close-btn" onclick={onClose}>&times;</button>
        </div>

        {#if isLoading}
            <div class="loading">Loading settings...</div>
        {:else}
            <form onsubmit={handleSubmit}>
                <div class="form-section">
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

                <div class="form-section">
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

                <div class="form-section">
                    <h3>Connection Settings</h3>

                    <div class="form-group checkbox-group">
                        <input
                                type="checkbox"
                                id="auto_start"
                                bind:checked={formData.connection.auto_start}
                        />
                        <label for="auto_start">Auto-start connection on launch</label>
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

                {#if error}
                    <div class="error-message">{error}</div>
                {/if}

                {#if saveSuccess}
                    <div class="success-message">Settings saved successfully!</div>
                {/if}

                <div class="form-actions">
                    <button type="button" class="btn btn-danger" onclick={handleRestartServer}>
                        Restart Server
                    </button>
                    <div class="spacer"></div>
                    <button type="button" class="btn btn-secondary" onclick={onClose}>
                        Cancel
                    </button>
                    <button
                            type="submit"
                            class="btn btn-primary"
                            disabled={isSaving}
                    >
                        {isSaving ? 'Saving...' : 'Save'}
                    </button>
                </div>
            </form>
        {/if}
    </div>
</div>

<style>
    .modal-backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(0, 0, 0, 0.5);
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 100;
    }

    .modal-content {
        background-color: #1b1e1f;
        border-radius: 4px;
        width: 90%;
        max-width: 600px;
        max-height: 90vh;
        overflow-y: auto;
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
    }

    .modal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem;
        border-bottom: 1px solid #262a2b;
        position: sticky;
        top: 0;
        background-color: #1b1e1f;
        z-index: 1;
    }

    .modal-header h2 {
        margin: 0;
        color: #aec2d3;
    }

    .close-btn {
        background: none;
        border: none;
        font-size: 1.5rem;
        color: #6c757d;
        cursor: pointer;
    }

    form {
        padding: 1rem;
    }

    .form-section {
        margin-bottom: 2rem;
    }

    .form-section h3 {
        margin-top: 0;
        margin-bottom: 1rem;
        color: #aec2d3;
        font-size: 1.1rem;
        border-bottom: 1px solid #262a2b;
        padding-bottom: 0.5rem;
    }

    .form-group {
        margin-bottom: 1rem;
    }

    label {
        display: block;
        margin-bottom: 0.5rem;
        color: #aec2d3;
    }

    input, select {
        width: 100%;
        padding: 0.5rem;
        border: 1px solid #262a2b;
        border-radius: 4px;
        background-color: #2b2a33;
        color: #fbfbfe;
    }

    .checkbox-group {
        display: flex;
        align-items: center;
    }

    .checkbox-group input {
        width: auto;
        margin-right: 0.5rem;
    }

    .form-actions {
        display: flex;
        justify-content: flex-end;
        gap: 0.5rem;
        margin-top: 1.5rem;
    }

    .spacer {
        flex-grow: 1;
    }

    .btn {
        padding: 0.5rem 1rem;
        border: none;
        border-radius: 4px;
        cursor: pointer;
        font-size: 1rem;
    }

    .btn:disabled {
        opacity: 0.7;
        cursor: not-allowed;
    }

    .btn-primary {
        background-color: #0062cc;
        color: white;
    }

    .btn-secondary {
        background-color: #6c757d;
        color: white;
    }

    .btn-danger {
        background-color: #dc3545;
        color: white;
    }

    .loading {
        padding: 2rem;
        text-align: center;
        color: #aec2d3;
    }

    .error-message {
        background-color: #482121;
        border: 1px solid #692929;
        color: #f5c2c2;
        padding: 0.75rem 1rem;
        border-radius: 4px;
        margin-top: 1rem;
    }

    .success-message {
        background-color: #1a3e29;
        border: 1px solid #245931;
        color: #b9ffc9;
        padding: 0.75rem 1rem;
        border-radius: 4px;
        margin-top: 1rem;
    }

    small {
        color: #6c757d;
        font-size: 0.875rem;
        display: block;
        margin-top: 0.25rem;
    }
</style>
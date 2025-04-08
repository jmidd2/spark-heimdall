<!-- frontend/src/lib/components/SettingsModal.svelte -->
<script lang="ts">
  import { api } from '$lib/api';
  import type { Device, AppConfig } from '$lib/types';
  import { onMount } from 'svelte';

  // Props
  type SettingsModalProps = {
    devices: Device[];
    onClose: () => void;
    showModal: boolean
  }
  let { devices, onClose, showModal = $bindable() }: SettingsModalProps = $props()

  let dialog = $state<HTMLDialogElement | undefined>();

  $effect(() => {
    if (showModal) dialog?.showModal();
  });

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

  function handleClose() {
    showModal = false;
    dialog?.close()
    onClose();
  }
</script>

<dialog class="modal"
        bind:this={dialog}
        onclose={() => (showModal = false)}
        onclick={(e) => { if (e.target === dialog) dialog.close(); }}>
    <div class="modal-content">
        <div class="modal-header">
            <h2>Application Settings</h2>
            <button type="button" class="close-btn" onclick={handleClose}>&times;</button>
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
                    <button type="button" class="btn btn-secondary" onclick={handleClose}>
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
</dialog>

<style>

</style>
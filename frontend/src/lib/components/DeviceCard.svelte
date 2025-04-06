<script lang="ts">
  import { api } from '$lib/api';
  import type { Device } from '$lib/types';

  export let device: Device;
  export let isConnected: boolean = false;
  export let onEdit: (device: Device) => void;
  export let onDelete: () => void;

  const handleConnect = async () => {
    try {
      if (isConnected) {
        await api.disconnect();
      } else {
        await api.connectToDevice(device.id);
      }
      // Reload page to reflect connection state
      window.location.reload();
    } catch (error) {
      alert(`Connection error: ${error instanceof Error ? error.message : String(error)}`);
    }
  };

  const handleDelete = async () => {
    if (confirm(`Are you sure you want to delete ${device.name}?`)) {
      try {
        await api.deleteDevice(device.id);
        onDelete();
      } catch (error) {
        alert(`Delete error: ${error instanceof Error ? error.message : String(error)}`);
      }
    }
  };
</script>

<div class="card {isConnected ? 'connected' : ''}">
    <div class="card-header">
        <h3 class="card-title">{device.name}</h3>
        <div class="card-actions">
            <button class="btn btn-secondary" on:click={() => onEdit(device)}>Edit</button>
            <button class="btn btn-danger" on:click={handleDelete}>Delete</button>
        </div>
    </div>
    <p>{device.ip_address}{device.port ? `:${device.port}` : ''} ({device.protocol})</p>
    {#if device.description}
        <p class="card-description">{device.description}</p>
    {/if}
    <button
            class="btn {isConnected ? 'btn-danger' : 'btn-primary'}"
            on:click={handleConnect}
    >
        {isConnected ? 'Disconnect' : 'Connect'}
    </button>
</div>

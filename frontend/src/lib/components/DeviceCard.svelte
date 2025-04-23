<script lang="ts">
import { api } from '$lib/api';
import type { Device } from '$lib/types';

type DeviceCardProps = {
  device: Device;
  isConnected: boolean;
  onConnect: (device: Device) => void;
  onEdit: (device: Device) => void;
  onDelete: () => Promise<void>;
};

const {
  device,
  isConnected = false,
  onDelete,
  onEdit,
  onConnect,
}: DeviceCardProps = $props();

const handleConnect = async () => {
  try {
    if (isConnected) {
      await api.disconnect();
    } else {
      await api.connectToDevice(device.id);
      onConnect(device);
    }
    // Reload page to reflect connection state
    window.location.reload();
  } catch (error) {
    alert(
      `Connection error: ${error instanceof Error ? error.message : String(error)}`
    );
  }
};

const handleDelete = async () => {
  await onDelete();
};
</script>

<div class="card {isConnected ? 'connected' : ''}">
    <div class="card-header">
        <h3 class="card-title">{device.name}</h3>
        <div class="card-header-actions">
            <button class="btn btn-secondary" onclick={() => onEdit(device)}>Edit</button>
            <button class="btn btn-danger" onclick={handleDelete}>Delete</button>
        </div>
    </div>
    <div class="card-content">
        <p>{device.ip_address}{device.port ? `:${ device.port }` : ''} <span class={['badge', device.protocol]}>{device.protocol}</span>
        </p>
        {#if device.description}
            <p class="card-description">{device.description}</p>
        {/if}
        <div class="card-actions">
            <button
                    class={['btn', {'btn-danger': isConnected}, {'btn-primary': !isConnected}]}
                    onclick={handleConnect}
            >
                {isConnected ? 'Disconnect' : 'Connect'}
            </button>
        </div>
    </div>
</div>

<style>
    .card {
        overflow: hidden;
        border: 1px solid var(--heimdall-border-color);
        border-radius: var(--heimdall-rounded);
        padding: var(--heimdall-spacing-lg);
        margin-bottom: var(--heimdall-spacing-xl);
        background-color: var(--heimdall-bg-card);
        display: grid;
        grid-template-columns: 1fr;

        &:hover {
            background-color: var(--heimdall-bg-card-hover);
        }
    }

    .card-header {
        display: grid;
        grid-template-columns: 2fr 1fr;
        justify-content: space-between;
        align-items: center;
        margin-bottom: var(--heimdall-spacing-sm);

        h3 {
            margin: 0;
            font-size: var(--heimdall-font-size-xl);
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .card-header-actions {
            display: flex;
            gap: var(--heimdall-spacing-sm);
            margin-left: var(--heimdall-spacing-sm);
            justify-content: flex-end;

            .btn {
                font-size: var(--heimdall-font-size-xs);
            }
        }
    }

    .card-content {
        display: flex;
        flex-direction: column;
        justify-content: space-between;

        .badge {
            border-radius: var(--heimdall-rounded-full);
            border: 1px solid;
            padding: calc(var(--heimdall-spacing-xs) / 2) var(--heimdall-spacing-sm);
            font-size: var(--heimdall-font-size-xs);

            &.vnc {
                border-color: var(--heimdall-badge-vnc-border);
                background-color: var(--heimdall-badge-vnc-bg);
            }

            &.rdp {
                border-color: var(--heimdall-badge-rdp-border);
                background-color: var(--heimdall-badge-rdp-bg);
            }
        }

        .card-description {
            color: var(--heimdall-text-muted);
            font-size: var(--heimdall-font-size-sm);
            margin-top: var(--heimdall-spacing-xs);
        }

        .card-actions {
            margin-top: var(--heimdall-spacing-md);
        }
    }

    .card:only-child {
        grid-column: 1 / 3;
    }
</style>

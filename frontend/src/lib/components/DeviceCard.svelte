<script lang="ts">
    import {api} from '$lib/api';
    import type {Device} from '$lib/types';

    // export let device: Device;
    // export let isConnected: boolean = false;
    // export let onEdit: (device: Device) => void;
    // export let onDelete: () => void;

    type DeviceCardProps = {
        device: Device;
        isConnected: boolean;
        onConnect: (device: Device) => void;
        onEdit: (device: Device) => void;
        onDelete: (deviceId: string) => void;
    }

    const {device, isConnected = false, onDelete, onEdit, onConnect}: DeviceCardProps = $props();

    const handleConnect = async () => {
        try {
            if (isConnected) {
                await api.disconnect();
            } else {
                await api.connectToDevice(device.id);
                onConnect(device)
            }
            // Reload page to reflect connection state
            window.location.reload();
        } catch (error) {
            alert(`Connection error: ${error instanceof Error ? error.message : String(error)}`);
        }
    };

    const handleDelete = async () => {
        onDelete();
        // if (confirm(`Are you sure you want to delete ${device.name}?`)) {
        //     try {
        //         await api.deleteDevice(device.id);
        //     } catch (error) {
        //         alert(`Delete error: ${error instanceof Error ? error.message : String(error)}`);
        //     }
        // }
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
        <p>{device.ip_address}{device.port ? `:${device.port}` : ''} <span class="badge" class:vnc={device.protocol === 'vnc'} class:rdp={device.protocol === 'rdp'} >{device.protocol}</span></p>
        {#if device.description}
            <p class="card-description">{device.description}</p>
        {/if}
        <div class="card-actions">
            <button
                    class="btn {isConnected ? 'btn-danger' : 'btn-primary'}"
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
        border: 1px solid #262a2b;
        border-radius: 4px;
        padding: 15px;
        margin-bottom: 20px;
        background-color: #1b1e1f;
        display: grid;
        grid-template-columns: 1fr;

        &:hover {
            background-color: #323539;
        }
    }

    .card-header {
        display: grid;
        grid-template-columns: 2fr 1fr;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 0.5rem;

        h3 {
            margin: 0;
            font-size: 1.2rem;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }

        .card-header-actions {
            display: flex;
            gap: 0.5rem;
            margin-left: 0.5rem;
            justify-content: flex-end;

            .btn {
                font-size: 0.75rem;
            }
        }
    }



    .card-content {
        display: flex;
        flex-direction: column;
        justify-content: space-between;

        .badge {
            border-radius: 9999px;
            border: 1px solid;
            padding: 0.125rem 0.5rem;
            font-size: 0.75rem;

            &.vnc {
                border-color: #1a3e29;
                background-color: #49a129;
            }

            &.rdp {
                border-color: #6c757d;
                background-color: #d9534f;
            }
        }

        .card-description {
            color: #939faa;
            font-size: 0.9rem;
            margin-top: 5px;

        }

        .card-actions {
            margin-top: 0.75rem;
        }
    }

    .card:only-child {
        grid-column: 1 / 3;
    }
</style>

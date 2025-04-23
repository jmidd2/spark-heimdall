<!-- frontend/src/lib/components/DeviceList.svelte -->
<script lang="ts">
import type { Device } from '$lib/types';
// biome-ignore lint/correctness/noUnusedImports: <explanation>
import DeviceCard from './DeviceCard.svelte';

// Props
// export let devices: Device[] = [];
// export let connectedId: string | null = null;
// export let onConnect: (device: Device) => void;
// export let onEdit: (device: Device) => void;
// export let onDelete: (deviceId: string) => void;
//
// // Optional grouping
// export let groupByProtocol: boolean = false;

type DeviceListProps = {
  devices: Device[];
  connectedId: string | null;
  onConnect: (device: Device) => void;
  onEdit: (device: Device) => void;
  onDelete: (deviceId: string) => void;
  groupByProtocol?: boolean;
};

const {
  devices = [],
  connectedId = null,
  onConnect,
  onEdit,
  onDelete,
  groupByProtocol,
}: DeviceListProps = $props();

// Computed properties
const sortedDevices = $derived(
  [...devices].sort((a, b) => a.name.localeCompare(b.name))
);

const deviceGroups = $derived(
  groupByProtocol
    ? {
        vnc: sortedDevices.filter(d => d.protocol === 'vnc'),
        rdp: sortedDevices.filter(d => d.protocol === 'rdp'),
      }
    : { all: sortedDevices }
);
</script>

<div class="device-list">
    {#if groupByProtocol}
        <!-- Grouped by protocol -->
        {#if deviceGroups.vnc && deviceGroups.vnc.length > 0}
            <div class="group">
                <h3 class="group-title">VNC Connections</h3>
                <div class="cards">
                    {#each deviceGroups.vnc as device (device.id)}
                        <DeviceCard
                                {device}
                                isConnected={device.id === connectedId}
                                onConnect={() => onConnect(device)}
                                onEdit={() => onEdit(device)}
                                onDelete={() => onDelete(device.id)}
                        />
                    {/each}
                </div>
            </div>
        {/if}

        {#if deviceGroups.rdp && deviceGroups.rdp.length > 0}
            <div class="group">
                <h3 class="group-title">RDP Connections</h3>
                <div class="cards">
                    {#each deviceGroups.rdp as device (device.id)}
                        <DeviceCard
                                {device}
                                isConnected={device.id === connectedId}
                                onConnect={() => onConnect(device)}
                                onEdit={() => onEdit(device)}
                                onDelete={() => onDelete(device.id)}
                        />
                    {/each}
                </div>
            </div>
        {/if}
    {:else}
        <!-- Flat list -->
        <div class="cards">
            {#each sortedDevices as device (device.id)}
                <DeviceCard
                        {device}
                        isConnected={device.id === connectedId}
                        onConnect={() => onConnect(device)}
                        onEdit={() => onEdit(device)}
                        onDelete={() => onDelete(device.id)}
                />
            {/each}
        </div>
    {/if}
</div>

<style>
    .device-list {
        margin: 1rem 0;
    }

    .group {
        margin-bottom: 2rem;
    }

    .group-title {
        color: #aec2d3;
        font-size: 1.2rem;
        margin-bottom: 1rem;
        padding-bottom: 0.5rem;
        border-bottom: 1px solid #262a2b;
    }

    .cards {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 1rem;
    }

    /* Responsive adjustments */
    @media (max-width: 768px) {
        .cards {
            grid-template-columns: 1fr;
        }
    }
</style>
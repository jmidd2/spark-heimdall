<!-- frontend/src/lib/components/DeviceList.svelte -->
<script lang="ts">
import type { Device } from '$lib/types';
// biome-ignore lint/correctness/noUnusedImports: <explanation>
import DeviceCard from './DeviceCard.svelte';

type DeviceListProps = {
  devices: Device[];
  connectedId: string | null;
  onConnect: (device: Device) => Promise<void>;
  onEdit: (device: Device) => Promise<void>;
  onDelete: (deviceId: string) => Promise<void>;
  groupByProtocol?: boolean;
};

let {
  devices = [],
  connectedId = null,
  onConnect,
  onEdit,
  onDelete,
  groupByProtocol = false,
}: DeviceListProps = $props();

let totalCount = $derived(devices.length);

// Computed properties
const sortedDevices = $derived(
  [...devices].sort((a, b) => a.name.localeCompare(b.name))
);

let filter: Device['protocol'] | 'none' = $state('none');

const filteredDevices = $derived(
  filter !== 'none'
    ? [...sortedDevices].filter(d => d.protocol === filter)
    : [...sortedDevices]
);

const deviceGroups = $derived(
  groupByProtocol
    ? {
        vnc: sortedDevices.filter(d => d.protocol === 'vnc'),
        rdp: sortedDevices.filter(d => d.protocol === 'rdp'),
      }
    : { all: sortedDevices }
);

function handleChange() {
  if (groupByProtocol) filter = 'none';
}
</script>

{#snippet deviceCard( d: Device )}
    <DeviceCard
            device={d}
            connectedId={connectedId}
            onConnect={onConnect}
            onEdit={onEdit}
            onDelete={onDelete}
    />
{/snippet}

{#snippet cardGroup(devices: Device[], title: string)}
    <div class="group">
        <h3 class="group-title">
            {title}
        </h3>
        <div class="cards">
            {#each devices as device (device.id)}
                {@render deviceCard( device )}
            {/each}
        </div>
    </div>
{/snippet}

<div class="device-list">
    <div class="filter-row">
        <div style="display: flex; align-items: center; gap: var(--heimdall-spacing-sm);)">
            <label aria-disabled={groupByProtocol} for="filterByProtocol">Filter:</label>
            <select bind:value={filter} id="filterByProtocol" disabled={groupByProtocol}>
                <option value="none">None</option>
                <option value="vnc">VNC</option>
                <option value="rdp">RDP</option>
            </select>
        </div>
        <div class="checkbox-group">
            <input type="checkbox" bind:checked={groupByProtocol} id="groupByProtocol" class="size-lg"
                   onchange={handleChange}>
            <label for="groupByProtocol" style="margin: 0">Group By Protocol</label>
        </div>
    </div>
    {#if groupByProtocol}
        <!-- Grouped by protocol -->
        {#if deviceGroups.vnc && deviceGroups.vnc.length > 0}
            {@render cardGroup(deviceGroups.vnc, 'VNC Connections')}
        {/if}

        {#if deviceGroups.rdp && deviceGroups.rdp.length > 0}
            {@render cardGroup(deviceGroups.rdp, 'RDP Connections')}
        {/if}
    {:else}
        <!-- Flat list -->
        {@render cardGroup(filteredDevices, filter === 'none' ? 'Connections' : `${filter.toUpperCase()} Connections`)}

        {#if ( filter !== 'none' )}
            <p>Showing {filteredDevices.length} of {totalCount} devices</p>
        {/if}
    {/if}
</div>

<style>
    .device-list {
        /*margin: var(--heimdall-spacing-lg) 0;*/

        .filter-row {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: var(--heimdall-spacing-sm);
            margin-bottom: var(--heimdall-spacing-lg);
            font-size: var(--heimdall-font-size-sm);
            border-bottom: 1px solid var(--heimdall-border-color);
            padding-bottom: var(--heimdall-spacing-sm);

            input, select {
                &[disabled] {
                    background-color: var(--heimdall-bg-cardg);
                    color: var(--heimdall-text-muted);
                }
            }

            label {
                margin: 0;

                &[aria-disabled="true"] {
                    color: var(--heimdall-text-muted);
                }
            }
        }
    }

    .group {
        margin-bottom: var(--heimdall-spacing-2xl);
    }

    .group-title {
        color: var(--heimdall-text-heading);
        font-size: var(--heimdall-font-size-xl);
        margin-bottom: var(--heimdall-spacing-lg);
        padding-bottom: var(--heimdall-spacing-sm);
        border-bottom: 1px solid var(--heimdall-border-color);
    }

    .cards {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: var(--heimdall-spacing-lg);
    }

    /* Responsive adjustments */
    @media (max-width: var(--heimdall-breakpoint-mobile)) {
        .cards {
            grid-template-columns: 1fr;
        }
    }
</style>
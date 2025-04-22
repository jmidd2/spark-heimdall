<!-- frontend/src/lib/components/DeviceModal.svelte -->
<script lang="ts">
    import type {Device} from '$lib/types';

    type DeviceModalProps = {
        device: Device | null;
        onSave: (device: Device) => void;
        onClose: () => void;
        showModal: boolean;
    };

    // Props
    let {
        device = null,
        onSave,
        onClose,
        showModal = $bindable(),
    }: DeviceModalProps = $props();

    // biome-ignore lint/style/useConst: <explanation>
    let dialog = $state<HTMLDialogElement | undefined>();

    $effect(() => {
        if (showModal) dialog?.showModal();
        if (!showModal) dialog?.close();
    });

    // Local state using Svelte 5 runes
    const formData = $state<Device>({
        id: '',
        name: '',
        ip_address: '',
        protocol: 'vnc',
        port: 0,
        username: '',
        password: '',
        full_screen: true,
        description: '',
        screen: '',
    });

    let errors = $state<Record<string, string>>({});

    // Initialize form data when component mounts or device prop changes
    $effect(() => {
        if (device) {
            // Editing existing device
            formData.id = device.id;
            formData.name = device.name;
            formData.ip_address = device.ip_address;
            formData.protocol = device.protocol;
            formData.port = device.port;
            formData.username = device.username || '';
            formData.password = ''; // Never display existing password
            formData.full_screen = device.full_screen;
            formData.description = device.description || '';
            formData.screen = device.screen || '';
        } else {
            // Adding new device
            formData.id = '';
            formData.name = '';
            formData.ip_address = '';
            formData.protocol = 'vnc';
            formData.port = 0;
            formData.username = '';
            formData.password = '';
            formData.full_screen = true;
            formData.description = '';
            formData.screen = '';
        }
    });

    function validateForm(): boolean {
        const newErrors: Record<string, string> = {};

        if (!formData.name.trim()) {
            newErrors.name = 'Name is required';
        }

        if (!formData.ip_address.trim()) {
            newErrors.ip_address = 'IP Address is required';
        }

        if (formData.protocol === 'rdp' && !formData.username.trim()) {
            newErrors.username = 'Username is required for RDP connections';
        }

        if (formData.port < 0) {
            newErrors.port = 'Port must be a positive number';
        }

        errors = newErrors;
        return Object.keys(newErrors).length === 0;
    }

    function handleSubmit(e: SubmitEvent) {
        e.preventDefault();

        if (validateForm()) {
            onSave({
                ...formData,
                // Convert port to number
                port:
                    typeof formData.port === 'string'
                        ? Number.parseInt(formData.port as string, 10) || 0
                        : formData.port,
            });
        }
    }

    function handleClose() {
        showModal = false;
        dialog?.close();
        onClose();
    }
</script>

<dialog bind:this={dialog} class="modal"
        onclose={() => (showModal = false)}
        onclick={(e) => { if (e.target === dialog) dialog.close(); }}>
    <div class="modal-content">
        <div class="modal-header">
            <h2>{device ? 'Edit Device' : 'Add New Device'}</h2>
            <button type="button" class="close-btn" onclick={handleClose}>&times;</button>
        </div>

        <form onsubmit={handleSubmit}>
            <div class="form-group">
                <label for="name">Name</label>
                <input
                        type="text"
                        id="name"
                        bind:value={formData.name}
                        placeholder="My Work PC"
                        class={errors.name ? 'error' : ''}
                />
                {#if errors.name}
                    <div class="error-message">{errors.name}</div>
                {/if}
            </div>

            <div class="form-group">
                <label for="ip_address">IP Address or Hostname</label>
                <input
                        type="text"
                        id="ip_address"
                        bind:value={formData.ip_address}
                        placeholder="192.168.1.100"
                        class={errors.ip_address ? 'error' : ''}
                />
                {#if errors.ip_address}
                    <div class="error-message">{errors.ip_address}</div>
                {/if}
            </div>

            <div class="form-group">
                <label for="protocol">Protocol</label>
                <select id="protocol" bind:value={formData.protocol}>
                    <option value="vnc">VNC</option>
                    <option value="rdp">RDP</option>
                </select>
            </div>

            <div class="form-group">
                <label for="port">Port (0 for default)</label>
                <input
                        type="number"
                        id="port"
                        bind:value={formData.port}
                        min="0"
                        max="65535"
                        class={errors.port ? 'error' : ''}
                />
                {#if errors.port}
                    <div class="error-message">{errors.port}</div>
                {/if}
            </div>

            <!-- Show username/password fields for RDP -->
            {#if formData.protocol === 'rdp'}
                <div class="form-group">
                    <label for="username">Username</label>
                    <input
                            type="text"
                            id="username"
                            bind:value={formData.username}
                            class={errors.username ? 'error' : ''}
                    />
                    {#if errors.username}
                        <div class="error-message">{errors.username}</div>
                    {/if}
                </div>

                <div class="form-group">
                    <label for="password">Password</label>
                    <input
                            type="password"
                            id="password"
                            bind:value={formData.password}
                            placeholder={device ? "••••••••" : ""}
                    />
                    {#if device}
                        <small>Leave blank to keep current password</small>
                    {/if}
                </div>
            {/if}

            <div class="form-group checkbox-group">
                <input type="checkbox" id="full_screen" bind:checked={formData.full_screen}/>
                <label for="full_screen">Launch in Full Screen</label>
            </div>

            <div class="form-group">
                <label for="description">Description (optional)</label>
                <input type="text" id="description" bind:value={formData.description}/>
            </div>

            {#if formData.protocol === 'vnc'}
                <div class="form-group">
                    <label for="screen">Screen Number (optional)</label>
                    <input type="text" id="screen" bind:value={formData.screen} placeholder="0"/>
                </div>
            {/if}

            <div class="form-actions">
                <div class="btn-group">
                    <button type="submit" class="btn btn-primary">Save</button>
                    <button type="button" class="btn btn-secondary" onclick={handleClose}>Cancel</button>
                </div>
            </div>
        </form>
    </div>
</dialog>

<style>
</style>
import { api } from '$lib/api';
import type { Device } from '$lib/types';

class DeviceStore {
  private _allDevices: Device[] = $state<Device[]>([]);
  private _filteredDevices: Device[] = $derived(this._allDevices);
  private _groupByProtocol = $state(false);
  private _protocolFilter: Device['protocol'] | 'none' = $state('none');
  private _totalCount = $derived(this._allDevices.length);
  private _filteredCount = $derived(this._filteredDevices.length);
  private _loadingState = $state(true);

  private _deviceGroups = $derived(
    this._groupByProtocol
      ? {
          vnc: this._allDevices.filter(d => d.protocol === 'vnc'),
          rdp: this._allDevices.filter(d => d.protocol === 'rdp'),
        }
      : { all: this._allDevices }
  );

  async initialize() {
    const result = await this.getDevices();
    this._allDevices = result.sort(this.sortByName);
    this._loadingState = false;
  }

  async getDevices() {
    return await api.getDevices();
  }

  get all() {
    return this._allDevices;
  }

  get filtered() {
    return this._filteredDevices;
  }

  get groupByProtocol() {
    return this._groupByProtocol;
  }

  set groupByProtocol(groupBy: boolean) {
    this._groupByProtocol = groupBy;
  }

  get protocolFilter() {
    return this._protocolFilter;
  }

  set protocolFilter(protocolFilter: Device['protocol'] | 'none') {
    this._protocolFilter = protocolFilter;
    if (this._protocolFilter === 'none') {
      this._filteredDevices = this._allDevices;
    } else {
      this.filterDevicesByKey('protocol', this._protocolFilter);
    }
  }

  get groups() {
    return this._deviceGroups;
  }

  get vnc() {
    return this._deviceGroups.vnc;
  }

  get rdp() {
    return this._deviceGroups.rdp;
  }

  get count() {
    return {
      total: this._totalCount,
      filtered: this._filteredCount,
      vnc: this.vnc?.length ?? 0,
      rdp: this.rdp?.length ?? 0,
    };
  }

  get isLoading() {
    return this._loadingState;
  }

  private sortByName(a: { name: string }, b: { name: string }) {
    return a.name.localeCompare(b.name);
  }

  private filterDevicesByKey<K extends keyof Device>(
    key: K,
    filter: Device[K] | string
  ): void {
    this._filteredDevices = this._allDevices
      .filter(d => d[key] === filter)
      .sort(this.sortByName);
  }

  async delete(deviceId: string) {
    this._loadingState = true;
    await api.deleteDevice(deviceId);
    this._allDevices = this._allDevices
      .filter(d => d.id !== deviceId)
      .sort(this.sortByName);
    this._loadingState = false;
  }

  async update(device: Device) {
    this._loadingState = true;
    const updated = await api.updateDevice(device);

    this._allDevices = this._allDevices
      .map(d => (d.id === device.id ? updated : d))
      .sort(this.sortByName);
    this._loadingState = false;
  }

  async add(device: Device) {
    this._loadingState = true;
    const added = await api.addDevice(device);
    this._allDevices = [...this._allDevices, added].sort(this.sortByName);
    this._loadingState = false;
  }

  search(id: string) {
    return this._allDevices.find(d => d.id === id);
  }
}

export const devices = new DeviceStore();
await devices.initialize();

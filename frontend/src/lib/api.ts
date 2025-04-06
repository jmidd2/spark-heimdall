// src/lib/api.ts
import type { ApiEndpoint, ApiResponse, AppConfig, Device } from './types';

function isApiError(json: ApiResponse<unknown>): json is { success: false, error: string } {
  return !( json as ApiResponse<unknown>
  ).success;
}

function hasData<T>(json: ApiResponse<T>): json is { success: true, data: T} {
  return (json as ApiResponse<unknown>).success && !!(json as ApiResponse<unknown>).data;
}

/**
 * API client for communicating with the Heimdall Go backend
 */
class ApiClient {
  /**
   * Base URL for API requests
   */
  private readonly baseUrl: string = '';

  constructor(baseUrl: string) {
    if (!baseUrl || baseUrl.length === 0) {
      throw new Error('Missing baseUrl');
    }
    this.baseUrl = baseUrl;
  }

  /**
   * Given a string that matches an ApiEndpoint build the url string for the fetch
   * @param {ApiEndpoint} url
   * @returns {string}
   * @private
   */
  private buildUrl(url: ApiEndpoint): string {
    let fetchUrl = this.baseUrl;

    if (url.startsWith('/connect')|| url === 'disconnect') {
      fetchUrl += url;
    } else {
      fetchUrl += `/api/${url}`;
    }

    return fetchUrl;
  }

  /**
   * Check given if given response was a success or has an error
   * @param {Response} response
   * @returns {Promise<T>}
   * @private
   */
  private async checkResponse<T>(response: Response) {
    const json: ApiResponse<T> = await response.json();

    if (isApiError(json) && !hasData<T>(json)) {
      throw new Error(json.error);
    }

    if (!hasData<T>(json)) {
      throw new Error('There was a problem fetching data.');
    }

    return json.data;
  }

  /**
   * make a GET fetch to the given ApiEndpoint
   * @param {ApiEndpoint} url
   * @param {string} errorText
   * @returns {Promise<T>}
   * @private
   */
  private async getFetch<T = unknown>(url: ApiEndpoint, errorText?: string): Promise<T> {
    // urls are /api/?? or /connect or /disconnect
    const fetchUrl = this.buildUrl(url);

    const response = await fetch(fetchUrl);
    if (!response.ok) {
      throw new Error(errorText ? errorText : 'There was a problem fetching data.');
    }
    return await this.checkResponse(response);
  }

  /**
   * make a POST fetch to the given ApiEndpoint
   * @param {ApiEndpoint} url
   * @param {Record<string, any>} body
   * @param {string} errorText
   * @returns {Promise<T>}
   * @private
   */
  private async postFetch<T = unknown>(url: ApiEndpoint, body?: Record<string, any>, errorText?: string): Promise<T> {
    const fetchUrl = this.buildUrl(url);

    const response = await fetch(fetchUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: body? JSON.stringify(body) : null,
    });

    if (!response.ok) {
      throw new Error(errorText ? errorText : 'There was a problem posting data.');
    }

    return await this.checkResponse(response)
  }

  private async deleteFetch(url: ApiEndpoint, body?: Record<string, any>, errorText?: string) {
    const fetchUrl = this.buildUrl(url);

    const response = await fetch(fetchUrl, {
      method: 'DELETE',
      body: body ? JSON.stringify(body) : null,
    });

    if (!response.ok) {
      throw new Error(errorText ? errorText : 'There was a problem with the delete request.');
    }
  }

  private async putFetch<T = unknown>(url: ApiEndpoint, body?: Record<string, any>, errorText?: string): Promise<T> {
    const fetchUrl = this.buildUrl(url);

    const response = await fetch(fetchUrl, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: body ? JSON.stringify(body) : null,
    });

    if (!response.ok) {
      throw new Error(errorText ? errorText : 'There was a problem with the put request.');
    }

    return await this.checkResponse(response)
  }

  /**
   * Fetches all devices from the server
   */
  async getDevices(): Promise<Device[]> {
    const response = await this.getFetch<Device[]>(`devices`, '`Failed to fetch devices')

    return response || [];
  }

  /**
   * Adds a new device
   */
  async addDevice(device: Omit<Device, 'id'>): Promise<Device> {
    return await this.postFetch<Device>( `devices`, device, '`Failed to add device' )
  }

  /**
   * Updates an existing device
   */
  async updateDevice(device: Device) {
    return await this.putFetch<Device>(`devices/${device.id}`, device, `Failed to update device`);
  }

  /**
   * Deletes a device
   */
  async deleteDevice(id: string): Promise<void> {
    await this.deleteFetch(`devices/${id}`, undefined, '`Failed to delete device');
  }

  /**
   * Fetches application configuration
   */
  async getConfig(): Promise<AppConfig> {
    return await this.getFetch<AppConfig>(`config`, 'Failed to fetch config')
  }

  /**
   * Updates application configuration
   */
  async updateConfig(config: Partial<AppConfig>): Promise<void> {
    await this.putFetch(`config`, config, 'Failed to update config');
  }

  /**
   * Connects to a remote device
   */
  async connectToDevice(id: string): Promise<void> {
    await this.postFetch(`connect/${id}`, undefined, '`Failed to connect to device');
  }

  /**
   * Disconnects from the current remote connection
   */
  async disconnect(): Promise<void> {
    await this.postFetch(`disconnect`, undefined, '`Failed to disconnect');
  }
}

// Create a singleton instance
export const api = new ApiClient(import.meta.env.VITE_API_URL);
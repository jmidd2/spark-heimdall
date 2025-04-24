// src/lib/api.ts
import type {
  ApiEndpoint,
  ApiResponse,
  ApiResponseNoData,
  AppConfig,
  Device,
} from './types';

function isApiError(
  json: ApiResponse<unknown>
): json is { success: false; error: string } {
  return !(json as ApiResponse<unknown>).success;
}

function hasData<T>(json: ApiResponse<T>): json is { success: true; data: T } {
  return (
    (json as ApiResponse<unknown>).success &&
    (json as ApiResponse<unknown>).data !== undefined
  );
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

    if (url.startsWith('/connect') || url === 'disconnect') {
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
  private async checkResponseWithData<T>(response: Response): Promise<T> {
    const json: ApiResponse<T> = await response.json();
    if (isApiError(json) && !hasData<T>(json)) {
      throw new Error(json.error);
    }

    if (!hasData<T>(json)) {
      throw new Error('There was a problem fetching data.');
    }

    return json.data;
  }

  private async checkResponseWithoutData(
    response: Response
  ): Promise<ApiResponseNoData> {
    const json: ApiResponseNoData = await response.json();

    if (isApiError(json)) {
      throw new Error(json.error);
    }

    return json;
  }

  /**
   * make a GET fetch to the given ApiEndpoint
   * @param {ApiEndpoint} url
   * @param {string} errorText
   * @returns {Promise<T>}
   * @private
   */
  private async getFetch<T = unknown>(
    url: ApiEndpoint,
    errorText?: string
  ): Promise<T> {
    // urls are /api/?? or /connect or /disconnect
    const fetchUrl = this.buildUrl(url);

    const response = await fetch(fetchUrl);
    if (!response.ok) {
      throw new Error(
        errorText ? errorText : 'There was a problem fetching data.'
      );
    }

    return await this.checkResponseWithData(response);
  }

  /**
   * make a POST fetch to the given ApiEndpoint
   * @param {ApiEndpoint} url
   * @param {Record<string, any>} body
   * @param {string} errorText
   * @returns {Promise<T>}
   * @private
   */
  private async postFetch<T = unknown>(
    url: ApiEndpoint,
    body?: Record<string, unknown>,
    errorText?: string
  ): Promise<T> {
    const fetchUrl = this.buildUrl(url);

    const response = await fetch(fetchUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: body ? JSON.stringify(body) : null,
    });

    if (!response.ok) {
      throw new Error(
        errorText ? errorText : 'There was a problem posting data.'
      );
    }

    return await this.checkResponseWithData(response);
  }

  private async deleteFetch(
    url: ApiEndpoint,
    body?: Record<string, unknown>,
    errorText?: string
  ) {
    const fetchUrl = this.buildUrl(url);

    const response = await fetch(fetchUrl, {
      method: 'DELETE',
      body: body ? JSON.stringify(body) : null,
    });

    if (!response.ok) {
      throw new Error(
        errorText ? errorText : 'There was a problem with the delete request.'
      );
    }
  }

  private async putFetch(
    url: ApiEndpoint,
    body?: Record<string, unknown>,
    errorText?: string
  ): Promise<Response> {
    const fetchUrl = this.buildUrl(url);

    const response = await fetch(fetchUrl, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: body ? JSON.stringify(body) : null,
    });

    if (!response.ok) {
      throw new Error(
        errorText ? errorText : 'There was a problem with the put request.'
      );
    }

    return response;
  }

  private async putFetchReturnsData<T = unknown>(
    url: ApiEndpoint,
    body?: Record<string, unknown>,
    errorText?: string
  ): Promise<T> {
    const response = await this.putFetch(url, body, errorText);

    return await this.checkResponseWithData(response);
  }

  private async putFetchReturnsNoData(
    url: ApiEndpoint,
    body?: Record<string, unknown>,
    errorText?: string
  ): Promise<ApiResponseNoData> {
    const response = await this.putFetch(url, body, errorText);

    return await this.checkResponseWithoutData(response);
  }

  /**
   * Fetches all devices from the server
   */
  async getDevices(): Promise<Device[]> {
    const response = await this.getFetch<Device[]>(
      'devices',
      'Failed to fetch devices'
    );

    return response || [];
  }

  /**
   * Adds a new device
   */
  async addDevice(device: Omit<Device, 'id'>): Promise<Device> {
    return await this.postFetch<Device>(
      'devices',
      device,
      'Failed to add device'
    );
  }

  /**
   * Updates an existing device
   */
  async updateDevice(device: Device) {
    return await this.putFetchReturnsData<Device>(
      `devices/${device.id}`,
      device,
      'Failed to update device'
    );
  }

  /**
   * Deletes a device
   */
  async deleteDevice(id: string): Promise<void> {
    await this.deleteFetch(
      `devices/${id}`,
      undefined,
      'Failed to delete device'
    );
  }

  /**
   * Fetches application configuration
   */
  async getConfig(): Promise<AppConfig> {
    return await this.getFetch<AppConfig>('config', 'Failed to fetch config');
  }

  /**
   * Updates application configuration
   */
  async updateConfig(config: Partial<AppConfig>): Promise<ApiResponseNoData> {
    return await this.putFetchReturnsNoData(
      'config',
      config,
      'Failed to update config'
    );
  }

  /**
   * Connects to a remote device
   */
  async connectToDevice(id: string): Promise<void> {
    await this.postFetch(
      `connect/${id}`,
      undefined,
      'Failed to connect to device'
    );
  }

  /**
   * Disconnects from the current remote connection
   */
  async disconnect(): Promise<void> {
    await this.postFetch('disconnect', undefined, 'Failed to disconnect');
  }
}

// Create a singleton instance
export const api = new ApiClient(import.meta.env.VITE_API_URL);

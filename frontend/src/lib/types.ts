// src/lib/types.ts

/**
 * Represents a remote device connection (PC, server, etc.)
 */
export type Device = {
  /** Unique identifier for the device */
  id: string;

  /** User-friendly name */
  name: string;

  /** IP address or hostname */
  ip_address: string;

  /** Connection protocol */
  protocol: 'vnc' | 'rdp';

  /** Connection port (0 for default) */
  port: number;

  /** Username for authentication (typically for RDP) */
  username?: string;

  /** Password for authentication (typically for RDP) */
  password?: string;

  /** Whether to launch in full screen mode */
  full_screen: boolean;

  /** Optional description */
  description?: string;

  /** Optional screen number for multi-monitor setups */
  screen?: string;
};

/**
 * Application configuration
 */
export type AppConfig = {
  /** Server configuration */
  server: {
    /** HTTP port for the web interface */
    port: number;
  };

  /** Connection settings */
  connection: {
    /** Whether to auto-start a connection when the app launches */
    auto_start: boolean;

    /** ID of the device to auto-connect to */
    auto_start_id: string;
  };

  /** External client applications */
  clients: {
    /** Path to VNC viewer executable */
    vnc_viewer: string;

    /** Path to VNC password file */
    vnc_password_file: string;

    /** Path to RDP client executable */
    rdp_viewer: string;
  };

  /** Logging configuration */
  logging: {
    /** Log level (debug, info, warn, error) */
    level: string;

    /** Log format (text, json) */
    format: string;
  };
};

/**
 * API response format
 */
export type ApiResponse<T> = {
  /** Whether the request was successful */
  success: boolean;

  /** Response data (if successful) */
  data?: T;

  /** Error message (if unsuccessful) */
  error?: string;
};

/**
 * API response format
 */
export type ApiResponseNoData = {
  /** Whether the request was successful */
  success: boolean;

  /** Error message (if unsuccessful) */
  error?: string;
};

/**
 * API Endpoints
 */
export type ApiEndpoint =
  | `connect`
  | `connect/${string}`
  | 'disconnect'
  | 'devices'
  | `devices/${string}`
  | 'config';

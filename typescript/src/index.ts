/**
 * Typed Environment Variables Library
 * 
 * A type-safe environment variable library that provides automatic type conversion
 * and validation for environment variables.
 * 
 * @author Parthiv Rawat
 * @license MIT
 */

/**
 * Base error class for environment variable errors
 */
export class EnvVarError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'EnvVarError';
  }
}

/**
 * Error thrown when required environment variable is not found
 */
export class EnvVarNotFoundError extends EnvVarError {
  constructor(key: string) {
    super(`Required environment variable '${key}' is not set`);
    this.name = 'EnvVarNotFoundError';
  }
}

/**
 * Error thrown when environment variable cannot be converted to expected type
 */
export class EnvVarTypeError extends EnvVarError {
  constructor(key: string, value: string, expectedType: string, details?: string) {
    const message = details
      ? `Environment variable '${key}' has value '${value}' which cannot be converted to ${expectedType}. ${details}`
      : `Environment variable '${key}' has value '${value}' which cannot be converted to ${expectedType}`;
    super(message);
    this.name = 'EnvVarTypeError';
  }
}

/**
 * Type-safe environment variable accessor
 */
export class Env {
  /**
   * Get raw environment variable value
   */
  private static getRaw(key: string, defaultValue?: string, required: boolean = true): string | undefined {
    const value = process.env[key];
    
    if (value === undefined) {
      if (required && defaultValue === undefined) {
        throw new EnvVarNotFoundError(key);
      }
      return defaultValue;
    }
    
    return value;
  }

  /**
   * Get string environment variable
   * 
   * @param key - Environment variable name
   * @param defaultValue - Default value if not set
   * @returns String value
   * @throws {EnvVarNotFoundError} If required variable is not set
   */
  static str(key: string): string;
  static str(key: string, defaultValue: string): string;
  static str(key: string, defaultValue?: string): string {
    const value = this.getRaw(key, defaultValue, defaultValue === undefined);
    return value ?? defaultValue!;
  }

  /**
   * Get integer environment variable
   * 
   * @param key - Environment variable name
   * @param defaultValue - Default value if not set
   * @returns Integer value
   * @throws {EnvVarNotFoundError} If required variable is not set
   * @throws {EnvVarTypeError} If value cannot be converted to int
   */
  static int(key: string): number;
  static int(key: string, defaultValue: number): number;
  static int(key: string, defaultValue?: number): number {
    const rawValue = this.getRaw(
      key,
      defaultValue !== undefined ? String(defaultValue) : undefined,
      defaultValue === undefined
    );
    
    if (rawValue === undefined) {
      return defaultValue!;
    }
    
    const parsed = parseInt(rawValue, 10);
    
    if (isNaN(parsed)) {
      throw new EnvVarTypeError(key, rawValue, 'int');
    }
    
    return parsed;
  }

  /**
   * Get float environment variable
   * 
   * @param key - Environment variable name
   * @param defaultValue - Default value if not set
   * @returns Float value
   * @throws {EnvVarNotFoundError} If required variable is not set
   * @throws {EnvVarTypeError} If value cannot be converted to float
   */
  static float(key: string): number;
  static float(key: string, defaultValue: number): number;
  static float(key: string, defaultValue?: number): number {
    const rawValue = this.getRaw(
      key,
      defaultValue !== undefined ? String(defaultValue) : undefined,
      defaultValue === undefined
    );
    
    if (rawValue === undefined) {
      return defaultValue!;
    }
    
    const parsed = parseFloat(rawValue);
    
    if (isNaN(parsed)) {
      throw new EnvVarTypeError(key, rawValue, 'float');
    }
    
    return parsed;
  }

  /**
   * Get boolean environment variable
   * 
   * Accepts: true/false, yes/no, 1/0, on/off (case-insensitive)
   * 
   * @param key - Environment variable name
   * @param defaultValue - Default value if not set
   * @returns Boolean value
   * @throws {EnvVarNotFoundError} If required variable is not set
   * @throws {EnvVarTypeError} If value cannot be converted to bool
   */
  static bool(key: string): boolean;
  static bool(key: string, defaultValue: boolean): boolean;
  static bool(key: string, defaultValue?: boolean): boolean {
    const rawValue = this.getRaw(
      key,
      defaultValue !== undefined ? String(defaultValue) : undefined,
      defaultValue === undefined
    );
    
    if (rawValue === undefined) {
      return defaultValue!;
    }
    
    const trueValues = new Set(['true', 'yes', '1', 'on', 't', 'y']);
    const falseValues = new Set(['false', 'no', '0', 'off', 'f', 'n']);
    
    const normalized = rawValue.toLowerCase().trim();
    
    if (trueValues.has(normalized)) {
      return true;
    } else if (falseValues.has(normalized)) {
      return false;
    } else {
      const validValues = [...trueValues, ...falseValues].join(', ');
      throw new EnvVarTypeError(key, rawValue, 'bool', `Valid values: ${validValues}`);
    }
  }

  /**
   * Get list environment variable
   * 
   * @param key - Environment variable name
   * @param options - Options including default value and separator
   * @returns List of strings
   * @throws {EnvVarNotFoundError} If required variable is not set
   */
  static list(key: string, options?: { default?: string[]; separator?: string }): string[] {
    const separator = options?.separator ?? ',';
    const defaultValue = options?.default;
    
    const rawValue = this.getRaw(
      key,
      defaultValue !== undefined ? defaultValue.join(separator) : undefined,
      defaultValue === undefined
    );
    
    if (rawValue === undefined) {
      return defaultValue!;
    }
    
    if (!rawValue.trim()) {
      return [];
    }
    
    return rawValue.split(separator).map(item => item.trim());
  }

  /**
   * Get dictionary environment variable
   * 
   * Format: KEY1=VALUE1,KEY2=VALUE2
   * 
   * @param key - Environment variable name
   * @param options - Options including default value and separators
   * @returns Dictionary
   * @throws {EnvVarNotFoundError} If required variable is not set
   * @throws {EnvVarTypeError} If value cannot be parsed as dict
   */
  static dict(
    key: string,
    options?: {
      default?: Record<string, string>;
      itemSeparator?: string;
      keyValueSeparator?: string;
    }
  ): Record<string, string> {
    const itemSeparator = options?.itemSeparator ?? ',';
    const keyValueSeparator = options?.keyValueSeparator ?? '=';
    const defaultValue = options?.default;
    
    const defaultStr = defaultValue
      ? Object.entries(defaultValue)
          .map(([k, v]) => `${k}${keyValueSeparator}${v}`)
          .join(itemSeparator)
      : undefined;
    
    const rawValue = this.getRaw(key, defaultStr, defaultValue === undefined);
    
    if (rawValue === undefined) {
      return defaultValue!;
    }
    
    if (!rawValue.trim()) {
      return {};
    }
    
    const result: Record<string, string> = {};
    
    try {
      for (const item of rawValue.split(itemSeparator)) {
        const trimmedItem = item.trim();
        if (!trimmedItem) continue;
        
        const [k, v] = trimmedItem.split(keyValueSeparator, 2);
        if (v === undefined) {
          throw new Error('Invalid format');
        }
        result[k.trim()] = v.trim();
      }
      return result;
    } catch (error) {
      throw new EnvVarTypeError(
        key,
        rawValue,
        'dict',
        `Expected format: KEY1${keyValueSeparator}VALUE1${itemSeparator}KEY2${keyValueSeparator}VALUE2`
      );
    }
  }

  /**
   * Get enum environment variable
   * 
   * @param key - Environment variable name
   * @param enumObj - Enum object
   * @param defaultValue - Default value if not set
   * @returns Enum value
   * @throws {EnvVarNotFoundError} If required variable is not set
   * @throws {EnvVarTypeError} If value is not a valid enum member
   */
  static enum<T extends Record<string, string | number>>(
    key: string,
    enumObj: T,
    defaultValue?: T[keyof T]
  ): T[keyof T] {
    const defaultStr = defaultValue !== undefined ? String(defaultValue) : undefined;
    const rawValue = this.getRaw(key, defaultStr, defaultValue === undefined);
    
    if (rawValue === undefined) {
      return defaultValue!;
    }
    
    const enumKeys = Object.keys(enumObj);
    
    const upperValue = rawValue.toUpperCase();
    const matchingKey = enumKeys.find(k => k.toUpperCase() === upperValue);
    
    if (matchingKey) {
      return enumObj[matchingKey] as T[keyof T];
    }
    
    throw new EnvVarTypeError(
      key,
      rawValue,
      'enum',
      `Valid values: ${enumKeys.join(', ')}`
    );
  }

  /**
   * Get URL environment variable with basic validation
   * 
   * @param key - Environment variable name
   * @param defaultValue - Default value if not set
   * @returns URL string
   * @throws {EnvVarNotFoundError} If required variable is not set
   * @throws {EnvVarTypeError} If value is not a valid URL
   */
  static url(key: string): string;
  static url(key: string, defaultValue: string): string;
  static url(key: string, defaultValue?: string): string {
    const value = this.str(key, defaultValue as any);
    
    if (!value) {
      return value;
    }
    
    try {
      new URL(value);
      return value;
    } catch {
      throw new EnvVarTypeError(
        key,
        value,
        'URL',
        'Value is not a valid URL'
      );
    }
  }

  /**
   * Get environment variable with custom converter
   * 
   * @param key - Environment variable name
   * @param converter - Function to convert string to desired type
   * @param defaultValue - Default value if not set
   * @returns Converted value
   * @throws {EnvVarNotFoundError} If required variable is not set
   * @throws {EnvVarTypeError} If conversion fails
   */
  static custom<T>(
    key: string,
    converter: (value: string) => T,
    defaultValue?: T
  ): T {
    const rawValue = this.getRaw(key, undefined, defaultValue === undefined);
    
    if (rawValue === undefined) {
      return defaultValue!;
    }
    
    try {
      return converter(rawValue);
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : String(error);
      throw new EnvVarTypeError(key, rawValue, 'custom', errorMessage);
    }
  }
}

/**
 * Singleton instance for convenient access
 */
export const env = Env;

/**
 * Default export
 */
export default env;

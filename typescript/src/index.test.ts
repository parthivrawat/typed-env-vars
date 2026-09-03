/**
 * Tests for Typed Environment Variables Library
 */

import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { env, EnvVarNotFoundError, EnvVarTypeError } from './index';

enum LogLevel {
  DEBUG = 'DEBUG',
  INFO = 'INFO',
  WARNING = 'WARNING',
  ERROR = 'ERROR',
}

describe('Env.str', () => {
  beforeEach(() => {
    delete process.env.TEST_VAR;
  });

  it('should return string value when present', () => {
    process.env.TEST_VAR = 'hello';
    expect(env.str('TEST_VAR')).toBe('hello');
  });

  it('should throw when required variable is missing', () => {
    expect(() => env.str('MISSING_VAR')).toThrow(EnvVarNotFoundError);
  });

  it('should return default value when variable is missing', () => {
    expect(env.str('MISSING_VAR', 'default')).toBe('default');
  });

  it('should handle empty string', () => {
    process.env.TEST_VAR = '';
    expect(env.str('TEST_VAR')).toBe('');
  });
});

describe('Env.int', () => {
  beforeEach(() => {
    delete process.env.TEST_VAR;
  });

  it('should parse valid integer', () => {
    process.env.TEST_VAR = '42';
    expect(env.int('TEST_VAR')).toBe(42);
  });

  it('should parse negative integer', () => {
    process.env.TEST_VAR = '-10';
    expect(env.int('TEST_VAR')).toBe(-10);
  });

  it('should throw on invalid integer', () => {
    process.env.TEST_VAR = 'not_a_number';
    expect(() => env.int('TEST_VAR')).toThrow(EnvVarTypeError);
  });

  it('should return default value', () => {
    expect(env.int('MISSING_VAR', 100)).toBe(100);
  });

  it('should throw when required variable is missing', () => {
    expect(() => env.int('MISSING_VAR')).toThrow(EnvVarNotFoundError);
  });
});

describe('Env.float', () => {
  beforeEach(() => {
    delete process.env.TEST_VAR;
  });

  it('should parse valid float', () => {
    process.env.TEST_VAR = '3.14';
    expect(env.float('TEST_VAR')).toBe(3.14);
  });

  it('should parse scientific notation', () => {
    process.env.TEST_VAR = '1.5e-10';
    expect(env.float('TEST_VAR')).toBe(1.5e-10);
  });

  it('should throw on invalid float', () => {
    process.env.TEST_VAR = 'not_a_float';
    expect(() => env.float('TEST_VAR')).toThrow(EnvVarTypeError);
  });

  it('should return default value', () => {
    expect(env.float('MISSING_VAR', 2.5)).toBe(2.5);
  });
});

describe('Env.bool', () => {
  beforeEach(() => {
    delete process.env.TEST_VAR;
  });

  it('should parse true values', () => {
    const trueValues = ['true', 'True', 'TRUE', 'yes', 'YES', '1', 'on', 'ON', 't', 'y'];
    
    for (const value of trueValues) {
      process.env.TEST_VAR = value;
      expect(env.bool('TEST_VAR')).toBe(true);
    }
  });

  it('should parse false values', () => {
    const falseValues = ['false', 'False', 'FALSE', 'no', 'NO', '0', 'off', 'OFF', 'f', 'n'];
    
    for (const value of falseValues) {
      process.env.TEST_VAR = value;
      expect(env.bool('TEST_VAR')).toBe(false);
    }
  });

  it('should throw on invalid boolean', () => {
    process.env.TEST_VAR = 'maybe';
    expect(() => env.bool('TEST_VAR')).toThrow(EnvVarTypeError);
  });

  it('should return default value', () => {
    expect(env.bool('MISSING_VAR', true)).toBe(true);
    expect(env.bool('MISSING_VAR', false)).toBe(false);
  });
});

describe('Env.list', () => {
  beforeEach(() => {
    delete process.env.TEST_VAR;
  });

  it('should parse comma-separated list', () => {
    process.env.TEST_VAR = 'a,b,c';
    expect(env.list('TEST_VAR')).toEqual(['a', 'b', 'c']);
  });

  it('should trim whitespace', () => {
    process.env.TEST_VAR = 'a, b , c';
    expect(env.list('TEST_VAR')).toEqual(['a', 'b', 'c']);
  });

  it('should use custom separator', () => {
    process.env.TEST_VAR = 'a:b:c';
    expect(env.list('TEST_VAR', { separator: ':' })).toEqual(['a', 'b', 'c']);
  });

  it('should return empty array for empty string', () => {
    process.env.TEST_VAR = '';
    expect(env.list('TEST_VAR')).toEqual([]);
  });

  it('should return default value', () => {
    expect(env.list('MISSING_VAR', { default: ['x', 'y'] })).toEqual(['x', 'y']);
  });

  it('should not serialize defaults through the separator', () => {
    expect(env.list('MISSING_VAR', { default: ['a,b'] })).toEqual(['a,b']);
  });
});

describe('Env.dict', () => {
  beforeEach(() => {
    delete process.env.TEST_VAR;
  });

  it('should parse dictionary', () => {
    process.env.TEST_VAR = 'key1=value1,key2=value2';
    expect(env.dict('TEST_VAR')).toEqual({
      key1: 'value1',
      key2: 'value2',
    });
  });

  it('should trim whitespace', () => {
    process.env.TEST_VAR = 'key1 = value1 , key2 = value2';
    expect(env.dict('TEST_VAR')).toEqual({
      key1: 'value1',
      key2: 'value2',
    });
  });

  it('should use custom separators', () => {
    process.env.TEST_VAR = 'key1:value1;key2:value2';
    expect(
      env.dict('TEST_VAR', {
        itemSeparator: ';',
        keyValueSeparator: ':',
      })
    ).toEqual({
      key1: 'value1',
      key2: 'value2',
    });
  });

  it('should return empty object for empty string', () => {
    process.env.TEST_VAR = '';
    expect(env.dict('TEST_VAR')).toEqual({});
  });

  it('should throw on invalid format', () => {
    process.env.TEST_VAR = 'invalid_format';
    expect(() => env.dict('TEST_VAR')).toThrow(EnvVarTypeError);
  });

  it('should return default value', () => {
    const defaultValue = { a: '1', b: '2' };
    expect(env.dict('MISSING_VAR', { default: defaultValue })).toEqual(defaultValue);
  });

  it('should not serialize defaults through the separator', () => {
    const defaultValue = { a: '1,2' };
    expect(env.dict('MISSING_VAR', { default: defaultValue })).toEqual(defaultValue);
  });
});

describe('Env.enum', () => {
  beforeEach(() => {
    delete process.env.LOG_LEVEL;
  });

  it('should parse valid enum value', () => {
    process.env.LOG_LEVEL = 'INFO';
    expect(env.enum('LOG_LEVEL', LogLevel)).toBe(LogLevel.INFO);
  });

  it('should be case-insensitive', () => {
    process.env.LOG_LEVEL = 'debug';
    expect(env.enum('LOG_LEVEL', LogLevel)).toBe(LogLevel.DEBUG);
  });

  it('should throw on invalid enum value', () => {
    process.env.LOG_LEVEL = 'INVALID';
    expect(() => env.enum('LOG_LEVEL', LogLevel)).toThrow(EnvVarTypeError);
  });

  it('should return default value', () => {
    expect(env.enum('MISSING_VAR', LogLevel, LogLevel.WARNING)).toBe(LogLevel.WARNING);
  });
});

describe('Env.url', () => {
  beforeEach(() => {
    delete process.env.TEST_VAR;
  });

  it('should accept http URL', () => {
    process.env.TEST_VAR = 'http://example.com';
    expect(env.url('TEST_VAR')).toBe('http://example.com');
  });

  it('should accept https URL', () => {
    process.env.TEST_VAR = 'https://example.com/path';
    expect(env.url('TEST_VAR')).toBe('https://example.com/path');
  });

  it('should accept websocket URL', () => {
    process.env.TEST_VAR = 'ws://example.com';
    expect(env.url('TEST_VAR')).toBe('ws://example.com');
  });

  it('should throw on invalid URL', () => {
    process.env.TEST_VAR = 'not-a-url';
    expect(() => env.url('TEST_VAR')).toThrow(EnvVarTypeError);
  });

  it('should return default value', () => {
    expect(env.url('MISSING_VAR', 'https://default.com')).toBe('https://default.com');
  });
});

describe('Env.custom', () => {
  beforeEach(() => {
    delete process.env.TEST_VAR;
  });

  it('should use custom converter', () => {
    process.env.TEST_VAR = '2024-01-15';
    
    const parseDate = (s: string) => {
      const [year, month, day] = s.split('-').map(Number);
      return new Date(year, month - 1, day);
    };
    
    const result = env.custom('TEST_VAR', parseDate);
    expect(result.getFullYear()).toBe(2024);
    expect(result.getMonth()).toBe(0);
    expect(result.getDate()).toBe(15);
  });

  it('should throw on conversion error', () => {
    process.env.TEST_VAR = 'invalid';
    
    const parseInt = (s: string) => {
      const num = parseInt(s, 10);
      if (isNaN(num)) throw new Error('Not a number');
      return num;
    };
    
    expect(() => env.custom('TEST_VAR', parseInt)).toThrow(EnvVarTypeError);
  });

  it('should return default value', () => {
    const parseInt = (s: string) => parseInt(s, 10);
    expect(env.custom('MISSING_VAR', parseInt, 42)).toBe(42);
  });
});

describe('Real-world scenarios', () => {
  beforeEach(() => {
    process.env.DATABASE_URL = 'postgresql://localhost:5432/mydb';
    process.env.DB_POOL_SIZE = '20';
    process.env.DB_TIMEOUT = '30.5';
    process.env.DB_SSL = 'true';
  });

  afterEach(() => {
    delete process.env.DATABASE_URL;
    delete process.env.DB_POOL_SIZE;
    delete process.env.DB_TIMEOUT;
    delete process.env.DB_SSL;
  });

  it('should parse database configuration', () => {
    const dbUrl = env.url('DATABASE_URL');
    const poolSize = env.int('DB_POOL_SIZE');
    const timeout = env.float('DB_TIMEOUT');
    const ssl = env.bool('DB_SSL');

    expect(dbUrl).toBe('postgresql://localhost:5432/mydb');
    expect(poolSize).toBe(20);
    expect(timeout).toBe(30.5);
    expect(ssl).toBe(true);
  });
});

describe('Application configuration', () => {
  beforeEach(() => {
    process.env.APP_NAME = 'MyApp';
    process.env.DEBUG = 'false';
    process.env.ALLOWED_HOSTS = 'localhost,127.0.0.1,example.com';
    process.env.LOG_LEVEL = 'INFO';
  });

  afterEach(() => {
    delete process.env.APP_NAME;
    delete process.env.DEBUG;
    delete process.env.ALLOWED_HOSTS;
    delete process.env.LOG_LEVEL;
  });

  it('should parse application configuration', () => {
    const appName = env.str('APP_NAME');
    const debug = env.bool('DEBUG');
    const allowedHosts = env.list('ALLOWED_HOSTS');
    const logLevel = env.enum('LOG_LEVEL', LogLevel);

    expect(appName).toBe('MyApp');
    expect(debug).toBe(false);
    expect(allowedHosts).toEqual(['localhost', '127.0.0.1', 'example.com']);
    expect(logLevel).toBe(LogLevel.INFO);
  });
});

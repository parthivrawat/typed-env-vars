"""
Typed Environment Variables Library

A type-safe environment variable library that provides automatic type conversion
and validation for environment variables.

Author: Parthiv Rawat
License: MIT
"""

import os
from typing import TypeVar, Optional, Callable, Any, List, Dict
from enum import Enum
from urllib.parse import urlparse


T = TypeVar('T')


class EnvVarError(Exception):
    """Exception raised for environment variable errors."""
    pass


class EnvVarNotFoundError(EnvVarError):
    """Exception raised when required environment variable is not found."""
    pass


class EnvVarTypeError(EnvVarError):
    """Exception raised when environment variable cannot be converted to expected type."""
    pass


class Env:
    """
    Type-safe environment variable accessor.
    
    Usage:
        from typed_env import env
        
        DATABASE_URL = env.str('DATABASE_URL')
        MAX_CONNECTIONS = env.int('MAX_CONNECTIONS', default=10)
        DEBUG = env.bool('DEBUG', default=False)
        ALLOWED_HOSTS = env.list('ALLOWED_HOSTS', default=['localhost'])
    """
    
    @staticmethod
    def _get_raw(key: str, default: Optional[str] = None, required: bool = True) -> Optional[str]:
        """Get raw environment variable value."""
        value = os.environ.get(key)
        
        if value is None:
            if required and default is None:
                raise EnvVarNotFoundError(
                    f"Required environment variable '{key}' is not set"
                )
            return default
        
        return value
    
    @staticmethod
    def str(key: str, default: Optional[str] = None) -> str:
        """
        Get string environment variable.
        
        Args:
            key: Environment variable name
            default: Default value if not set
            
        Returns:
            String value
            
        Raises:
            EnvVarNotFoundError: If required variable is not set
        """
        value = Env._get_raw(key, default, required=(default is None))
        return value if value is not None else default
    
    @staticmethod
    def int(key: str, default: Optional[int] = None) -> int:
        """
        Get integer environment variable.
        
        Args:
            key: Environment variable name
            default: Default value if not set
            
        Returns:
            Integer value
            
        Raises:
            EnvVarNotFoundError: If required variable is not set
            EnvVarTypeError: If value cannot be converted to int
        """
        raw_value = Env._get_raw(key, str(default) if default is not None else None, 
                                  required=(default is None))
        
        if raw_value is None:
            return default
        
        try:
            return int(raw_value)
        except ValueError:
            raise EnvVarTypeError(
                f"Environment variable '{key}' has value '{raw_value}' "
                f"which cannot be converted to int"
            )
    
    @staticmethod
    def float(key: str, default: Optional[float] = None) -> float:
        """
        Get float environment variable.
        
        Args:
            key: Environment variable name
            default: Default value if not set
            
        Returns:
            Float value
            
        Raises:
            EnvVarNotFoundError: If required variable is not set
            EnvVarTypeError: If value cannot be converted to float
        """
        raw_value = Env._get_raw(key, str(default) if default is not None else None,
                                  required=(default is None))
        
        if raw_value is None:
            return default
        
        try:
            return float(raw_value)
        except ValueError:
            raise EnvVarTypeError(
                f"Environment variable '{key}' has value '{raw_value}' "
                f"which cannot be converted to float"
            )
    
    @staticmethod
    def bool(key: str, default: Optional[bool] = None) -> bool:
        """
        Get boolean environment variable.
        
        Accepts: true/false, yes/no, 1/0, on/off (case-insensitive)
        
        Args:
            key: Environment variable name
            default: Default value if not set
            
        Returns:
            Boolean value
            
        Raises:
            EnvVarNotFoundError: If required variable is not set
            EnvVarTypeError: If value cannot be converted to bool
        """
        raw_value = Env._get_raw(key, str(default) if default is not None else None,
                                  required=(default is None))
        
        if raw_value is None:
            return default
        
        true_values = {'true', 'yes', '1', 'on', 't', 'y'}
        false_values = {'false', 'no', '0', 'off', 'f', 'n'}
        
        normalized = raw_value.lower().strip()
        
        if normalized in true_values:
            return True
        elif normalized in false_values:
            return False
        else:
            raise EnvVarTypeError(
                f"Environment variable '{key}' has value '{raw_value}' "
                f"which cannot be converted to bool. "
                f"Valid values: {true_values | false_values}"
            )
    
    @staticmethod
    def list(key: str, default: Optional[List[str]] = None, 
             separator: str = ',') -> List[str]:
        """
        Get list environment variable.
        
        Args:
            key: Environment variable name
            default: Default value if not set
            separator: List item separator (default: ',')
            
        Returns:
            List of strings
            
        Raises:
            EnvVarNotFoundError: If required variable is not set
        """
        value = os.environ.get(key)
        if value is None:
            if default is None:
                raise EnvVarNotFoundError(
                    f"Required environment variable '{key}' is not set"
                )
            return list(default)
        
        if not value.strip():
            return []
        
        return [item.strip() for item in value.split(separator)]
    
    @staticmethod
    def dict(key: str, default: Optional[Dict[str, str]] = None,
             item_separator: str = ',', key_value_separator: str = '=') -> Dict[str, str]:
        """
        Get dictionary environment variable.
        
        Format: KEY1=VALUE1,KEY2=VALUE2
        
        Args:
            key: Environment variable name
            default: Default value if not set
            item_separator: Separator between key-value pairs (default: ',')
            key_value_separator: Separator between key and value (default: '=')
            
        Returns:
            Dictionary
            
        Raises:
            EnvVarNotFoundError: If required variable is not set
            EnvVarTypeError: If value cannot be parsed as dict
        """
        value = os.environ.get(key)
        if value is None:
            if default is None:
                raise EnvVarNotFoundError(
                    f"Required environment variable '{key}' is not set"
                )
            return dict(default)
        
        if not value.strip():
            return {}
        
        result = {}
        try:
            for item in value.split(item_separator):
                item = item.strip()
                if not item:
                    continue
                k, v = item.split(key_value_separator, 1)
                result[k.strip()] = v.strip()
            return result
        except ValueError:
            raise EnvVarTypeError(
                f"Environment variable '{key}' has value '{value}' "
                f"which cannot be parsed as dict. "
                f"Expected format: KEY1{key_value_separator}VALUE1{item_separator}KEY2{key_value_separator}VALUE2"
            )
    
    @staticmethod
    def enum(key: str, enum_class: type, default: Optional[Enum] = None) -> Enum:
        """
        Get enum environment variable.
        
        Args:
            key: Environment variable name
            enum_class: Enum class
            default: Default value if not set
            
        Returns:
            Enum value
            
        Raises:
            EnvVarNotFoundError: If required variable is not set
            EnvVarTypeError: If value is not a valid enum member
        """
        raw_value = Env._get_raw(key, 
                                  default.name if default is not None else None,
                                  required=(default is None))
        
        if raw_value is None:
            return default
        
        try:
            return enum_class[raw_value.upper()]
        except KeyError:
            valid_values = [e.name for e in enum_class]
            raise EnvVarTypeError(
                f"Environment variable '{key}' has value '{raw_value}' "
                f"which is not a valid {enum_class.__name__}. "
                f"Valid values: {valid_values}"
            )
    
    @staticmethod
    def url(key: str, default: Optional[str] = None) -> str:
        """
        Get URL environment variable with basic validation.
        
        Args:
            key: Environment variable name
            default: Default value if not set
            
        Returns:
            URL string
            
        Raises:
            EnvVarNotFoundError: If required variable is not set
            EnvVarTypeError: If value is not a valid URL
        """
        value = Env.str(key, default)
        
        if not value:
            return value
        
        parsed = urlparse(value)
        if not parsed.scheme or not parsed.netloc:
            raise EnvVarTypeError(
                f"Environment variable '{key}' has value '{value}' "
                f"which does not appear to be a valid URL"
            )
        
        return value
    
    @staticmethod
    def custom(key: str, converter: Callable[[str], T], 
               default: Optional[T] = None) -> T:
        """
        Get environment variable with custom converter.
        
        Args:
            key: Environment variable name
            converter: Function to convert string to desired type
            default: Default value if not set
            
        Returns:
            Converted value
            
        Raises:
            EnvVarNotFoundError: If required variable is not set
            EnvVarTypeError: If conversion fails
        """
        raw_value = Env._get_raw(key, None, required=(default is None))
        
        if raw_value is None:
            return default
        
        try:
            return converter(raw_value)
        except Exception as e:
            raise EnvVarTypeError(
                f"Environment variable '{key}' has value '{raw_value}' "
                f"which cannot be converted: {str(e)}"
            )


env = Env()

"""
Tests for Typed Environment Variables Library
"""

import os
import pytest
from enum import Enum
from typed_env import env, EnvVarNotFoundError, EnvVarTypeError


class LogLevel(Enum):
    DEBUG = 1
    INFO = 2
    WARNING = 3
    ERROR = 4


class TestStringEnv:
    def test_str_required_present(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'hello')
        assert env.str('TEST_VAR') == 'hello'
    
    def test_str_required_missing(self):
        with pytest.raises(EnvVarNotFoundError):
            env.str('MISSING_VAR')
    
    def test_str_with_default(self):
        assert env.str('MISSING_VAR', default='default') == 'default'
    
    def test_str_empty_string(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', '')
        assert env.str('TEST_VAR') == ''


class TestIntEnv:
    def test_int_valid(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', '42')
        assert env.int('TEST_VAR') == 42
    
    def test_int_negative(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', '-10')
        assert env.int('TEST_VAR') == -10
    
    def test_int_invalid(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'not_a_number')
        with pytest.raises(EnvVarTypeError):
            env.int('TEST_VAR')
    
    def test_int_with_default(self):
        assert env.int('MISSING_VAR', default=100) == 100
    
    def test_int_required_missing(self):
        with pytest.raises(EnvVarNotFoundError):
            env.int('MISSING_VAR')


class TestFloatEnv:
    def test_float_valid(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', '3.14')
        assert env.float('TEST_VAR') == 3.14
    
    def test_float_scientific(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', '1.5e-10')
        assert env.float('TEST_VAR') == 1.5e-10
    
    def test_float_invalid(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'not_a_float')
        with pytest.raises(EnvVarTypeError):
            env.float('TEST_VAR')
    
    def test_float_with_default(self):
        assert env.float('MISSING_VAR', default=2.5) == 2.5


class TestBoolEnv:
    def test_bool_true_values(self, monkeypatch):
        true_values = ['true', 'True', 'TRUE', 'yes', 'YES', '1', 'on', 'ON', 't', 'y']
        for value in true_values:
            monkeypatch.setenv('TEST_VAR', value)
            assert env.bool('TEST_VAR') is True, f"Failed for value: {value}"
    
    def test_bool_false_values(self, monkeypatch):
        false_values = ['false', 'False', 'FALSE', 'no', 'NO', '0', 'off', 'OFF', 'f', 'n']
        for value in false_values:
            monkeypatch.setenv('TEST_VAR', value)
            assert env.bool('TEST_VAR') is False, f"Failed for value: {value}"
    
    def test_bool_invalid(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'maybe')
        with pytest.raises(EnvVarTypeError):
            env.bool('TEST_VAR')
    
    def test_bool_with_default(self):
        assert env.bool('MISSING_VAR', default=True) is True
        assert env.bool('MISSING_VAR', default=False) is False


class TestListEnv:
    def test_list_comma_separated(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'a,b,c')
        assert env.list('TEST_VAR') == ['a', 'b', 'c']
    
    def test_list_with_spaces(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'a, b , c')
        assert env.list('TEST_VAR') == ['a', 'b', 'c']
    
    def test_list_custom_separator(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'a:b:c')
        assert env.list('TEST_VAR', separator=':') == ['a', 'b', 'c']
    
    def test_list_empty(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', '')
        assert env.list('TEST_VAR') == []
    
    def test_list_with_default(self):
        assert env.list('MISSING_VAR', default=['x', 'y']) == ['x', 'y']


class TestDictEnv:
    def test_dict_valid(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'key1=value1,key2=value2')
        result = env.dict('TEST_VAR')
        assert result == {'key1': 'value1', 'key2': 'value2'}
    
    def test_dict_with_spaces(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'key1 = value1 , key2 = value2')
        result = env.dict('TEST_VAR')
        assert result == {'key1': 'value1', 'key2': 'value2'}
    
    def test_dict_custom_separators(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'key1:value1;key2:value2')
        result = env.dict('TEST_VAR', item_separator=';', key_value_separator=':')
        assert result == {'key1': 'value1', 'key2': 'value2'}
    
    def test_dict_empty(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', '')
        assert env.dict('TEST_VAR') == {}
    
    def test_dict_invalid(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'invalid_format')
        with pytest.raises(EnvVarTypeError):
            env.dict('TEST_VAR')
    
    def test_dict_with_default(self):
        default = {'a': '1', 'b': '2'}
        assert env.dict('MISSING_VAR', default=default) == default


class TestEnumEnv:
    def test_enum_valid(self, monkeypatch):
        monkeypatch.setenv('LOG_LEVEL', 'INFO')
        assert env.enum('LOG_LEVEL', LogLevel) == LogLevel.INFO
    
    def test_enum_case_insensitive(self, monkeypatch):
        monkeypatch.setenv('LOG_LEVEL', 'debug')
        assert env.enum('LOG_LEVEL', LogLevel) == LogLevel.DEBUG
    
    def test_enum_invalid(self, monkeypatch):
        monkeypatch.setenv('LOG_LEVEL', 'INVALID')
        with pytest.raises(EnvVarTypeError):
            env.enum('LOG_LEVEL', LogLevel)
    
    def test_enum_with_default(self):
        assert env.enum('MISSING_VAR', LogLevel, default=LogLevel.WARNING) == LogLevel.WARNING


class TestUrlEnv:
    def test_url_http(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'http://example.com')
        assert env.url('TEST_VAR') == 'http://example.com'
    
    def test_url_https(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'https://example.com/path')
        assert env.url('TEST_VAR') == 'https://example.com/path'
    
    def test_url_websocket(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'ws://example.com')
        assert env.url('TEST_VAR') == 'ws://example.com'
    
    def test_url_invalid(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'not-a-url')
        with pytest.raises(EnvVarTypeError):
            env.url('TEST_VAR')
    
    def test_url_with_default(self):
        assert env.url('MISSING_VAR', default='https://default.com') == 'https://default.com'


class TestCustomEnv:
    def test_custom_converter(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', '2024-01-15')
        from datetime import datetime
        
        def parse_date(s):
            return datetime.strptime(s, '%Y-%m-%d').date()
        
        result = env.custom('TEST_VAR', parse_date)
        assert result.year == 2024
        assert result.month == 1
        assert result.day == 15
    
    def test_custom_converter_invalid(self, monkeypatch):
        monkeypatch.setenv('TEST_VAR', 'invalid')
        
        def parse_int(s):
            return int(s)
        
        with pytest.raises(EnvVarTypeError):
            env.custom('TEST_VAR', parse_int)
    
    def test_custom_with_default(self):
        def parse_int(s):
            return int(s)
        
        assert env.custom('MISSING_VAR', parse_int, default=42) == 42


class TestRealWorldScenarios:
    def test_database_config(self, monkeypatch):
        monkeypatch.setenv('DATABASE_URL', 'postgresql://localhost:5432/mydb')
        monkeypatch.setenv('DB_POOL_SIZE', '20')
        monkeypatch.setenv('DB_TIMEOUT', '30.5')
        monkeypatch.setenv('DB_SSL', 'true')
        
        db_url = env.url('DATABASE_URL')
        pool_size = env.int('DB_POOL_SIZE')
        timeout = env.float('DB_TIMEOUT')
        ssl = env.bool('DB_SSL')
        
        assert db_url == 'postgresql://localhost:5432/mydb'
        assert pool_size == 20
        assert timeout == 30.5
        assert ssl is True
    
    def test_app_config(self, monkeypatch):
        monkeypatch.setenv('APP_NAME', 'MyApp')
        monkeypatch.setenv('DEBUG', 'false')
        monkeypatch.setenv('ALLOWED_HOSTS', 'localhost,127.0.0.1,example.com')
        monkeypatch.setenv('LOG_LEVEL', 'INFO')
        
        app_name = env.str('APP_NAME')
        debug = env.bool('DEBUG')
        allowed_hosts = env.list('ALLOWED_HOSTS')
        log_level = env.enum('LOG_LEVEL', LogLevel)
        
        assert app_name == 'MyApp'
        assert debug is False
        assert allowed_hosts == ['localhost', '127.0.0.1', 'example.com']
        assert log_level == LogLevel.INFO
    
    def test_feature_flags(self, monkeypatch):
        monkeypatch.setenv('FEATURES', 'feature1=true,feature2=false,feature3=true')
        
        features = env.dict('FEATURES')
        assert features == {
            'feature1': 'true',
            'feature2': 'false',
            'feature3': 'true'
        }


if __name__ == '__main__':
    pytest.main([__file__, '-v'])

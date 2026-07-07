from python_backend.config import AppConfig
from python_backend.http_app import run_server


if __name__ == "__main__":
    run_server(AppConfig.from_env())

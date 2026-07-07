from __future__ import annotations

import json
from http import HTTPStatus
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from typing import Callable
from urllib.parse import urlparse

from .config import AppConfig
from .database import DatabaseInitializer, SQLiteConnectionFactory
from .repository import SQLiteRecordRepository
from .service import RecordService


class DashboardApplication:
    def __init__(self, config: AppConfig, record_service: RecordService) -> None:
        self._config = config
        self._record_service = record_service

    def make_handler(self) -> type[BaseHTTPRequestHandler]:
        app = self

        class DashboardRequestHandler(BaseHTTPRequestHandler):
            def do_GET(self) -> None:
                app.handle_get(self)

            def log_message(self, format: str, *args: object) -> None:
                print("%s - - %s" % (self.address_string(), format % args))

        return DashboardRequestHandler

    def handle_get(self, handler: BaseHTTPRequestHandler) -> None:
        path = urlparse(handler.path).path
        routes: dict[str, Callable[[BaseHTTPRequestHandler], None]] = {
            "/": self._serve_index,
            "/index.html": self._serve_index,
            "/api/records": self._serve_records,
            "/healthz": self._serve_health,
        }
        route = routes.get(path)
        if route is None:
            self._send_json(handler, {"error": "not found"}, HTTPStatus.NOT_FOUND)
            return
        route(handler)

    def _serve_index(self, handler: BaseHTTPRequestHandler) -> None:
        index_path = self._config.static_index
        if not index_path.exists():
            self._send_json(handler, {"error": "index.html not found"}, HTTPStatus.NOT_FOUND)
            return

        body = index_path.read_bytes()
        handler.send_response(HTTPStatus.OK)
        handler.send_header("Content-Type", "text/html; charset=utf-8")
        handler.send_header("Cache-Control", "no-store")
        handler.send_header("Content-Length", str(len(body)))
        handler.end_headers()
        handler.wfile.write(body)

    def _serve_records(self, handler: BaseHTTPRequestHandler) -> None:
        self._send_json(handler, self._record_service.list_records_for_api())

    def _serve_health(self, handler: BaseHTTPRequestHandler) -> None:
        self._send_json(handler, {"status": "ok"})

    def _send_json(
        self,
        handler: BaseHTTPRequestHandler,
        payload: object,
        status: HTTPStatus = HTTPStatus.OK,
    ) -> None:
        body = json.dumps(payload, ensure_ascii=False).encode("utf-8")
        handler.send_response(status)
        handler.send_header("Content-Type", "application/json; charset=utf-8")
        handler.send_header("Content-Length", str(len(body)))
        handler.end_headers()
        handler.wfile.write(body)


class ApplicationFactory:
    def __init__(self, config: AppConfig) -> None:
        self._config = config

    def create(self) -> DashboardApplication:
        connection_factory = SQLiteConnectionFactory(self._config.sqlite_path)
        DatabaseInitializer(connection_factory).initialize()
        repository = SQLiteRecordRepository(connection_factory)
        record_service = RecordService(repository)
        return DashboardApplication(self._config, record_service)


def run_server(config: AppConfig | None = None) -> None:
    app_config = config or AppConfig.from_env()
    app = ApplicationFactory(app_config).create()
    server = ThreadingHTTPServer(("", app_config.port), app.make_handler())
    print(f"Python AI R&D dashboard listening on http://localhost:{app_config.port}")
    print(f"SQLite database: {Path(app_config.sqlite_path)}")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\nPython AI R&D dashboard stopped.")
    finally:
        server.server_close()

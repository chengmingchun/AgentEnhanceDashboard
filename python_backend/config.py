from __future__ import annotations

import os
from dataclasses import dataclass
from pathlib import Path


@dataclass(frozen=True)
class AppConfig:
    port: int
    sqlite_path: Path
    static_index: Path

    @classmethod
    def from_env(cls) -> "AppConfig":
        project_root = Path(__file__).resolve().parents[1]
        sqlite_path = Path(os.getenv("SQLITE_PATH", project_root / "dashboard.db"))
        if not sqlite_path.is_absolute():
            sqlite_path = project_root / sqlite_path

        return cls(
            port=int(os.getenv("PORT", "8080")),
            sqlite_path=sqlite_path,
            static_index=project_root / "index.html",
        )

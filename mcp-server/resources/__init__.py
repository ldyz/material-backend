"""
MCP Server Resources Package
"""

from .data_providers import (
    ResourceProvider,
    InventoryAlertsProvider,
    PendingTasksProvider,
    MaterialPlansProvider,
    StockSummaryProvider,
    get_provider,
)

__all__ = [
    "ResourceProvider",
    "InventoryAlertsProvider",
    "PendingTasksProvider",
    "MaterialPlansProvider",
    "StockSummaryProvider",
    "get_provider",
]

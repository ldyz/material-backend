"""
Resource Data Providers for MCP Server
Provides real-time data resources for AI Agents
"""

import asyncio
from typing import Any
import httpx


class ResourceProvider:
    """Base class for resource providers"""

    def __init__(self, backend_url: str, headers: dict):
        self.backend_url = backend_url
        self.headers = headers
        self._cache = {}
        self._cache_time = {}

    async def get(self) -> dict[str, Any]:
        """Get resource data"""
        raise NotImplementedError


class InventoryAlertsProvider(ResourceProvider):
    """Provider for inventory alerts"""

    async def get(self) -> dict[str, Any]:
        """Get inventory alerts"""
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.backend_url}/api/agent/operate",
                headers=self.headers,
                json={
                    "operation": "query",
                    "resource": "stock",
                    "parameters": {"low_stock_alert": True},
                    "reasoning": "Resource: inventory://alerts",
                },
            )
            response.raise_for_status()
            return response.json()


class PendingTasksProvider(ResourceProvider):
    """Provider for pending workflow tasks"""

    async def get(self) -> dict[str, Any]:
        """Get pending workflow tasks"""
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.backend_url}/api/agent/operate",
                headers=self.headers,
                json={
                    "operation": "query",
                    "resource": "workflow",
                    "parameters": {"status": "pending"},
                    "reasoning": "Resource: workflow://pending-tasks",
                },
            )
            response.raise_for_status()
            return response.json()


class MaterialPlansProvider(ResourceProvider):
    """Provider for material plans overview"""

    async def get(self) -> dict[str, Any]:
        """Get material plans overview"""
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.backend_url}/api/agent/operate",
                headers=self.headers,
                json={
                    "operation": "query",
                    "resource": "material_plan",
                    "parameters": {},
                    "reasoning": "Resource: material://plans",
                },
            )
            response.raise_for_status()
            return response.json()


class StockSummaryProvider(ResourceProvider):
    """Provider for stock summary statistics"""

    async def get(self) -> dict[str, Any]:
        """Get stock summary"""
        async with httpx.AsyncClient() as client:
            response = await client.post(
                f"{self.backend_url}/api/agent/operate",
                headers=self.headers,
                json={
                    "operation": "analyze",
                    "resource": "inventory",
                    "parameters": {"question": "summary"},
                    "reasoning": "Resource: stock://summary",
                },
            )
            response.raise_for_status()
            return response.json()


def get_provider(
    resource_uri: str,
    backend_url: str,
    headers: dict,
) -> ResourceProvider:
    """Get resource provider by URI"""
    providers = {
        "inventory://alerts": InventoryAlertsProvider,
        "workflow://pending-tasks": PendingTasksProvider,
        "material://plans": MaterialPlansProvider,
        "stock://summary": StockSummaryProvider,
    }

    provider_class = providers.get(resource_uri)
    if provider_class:
        return provider_class(backend_url, headers)
    raise ValueError(f"Unknown resource URI: {resource_uri}")

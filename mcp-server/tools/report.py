"""
Report Tools for MCP Server
Provides tools for generating various reports
"""

from typing import Any
import httpx


async def generate_inventory_summary(
    backend_url: str,
    headers: dict,
) -> dict[str, Any]:
    """Generate inventory summary report

    Args:
        backend_url: Backend API base URL
        headers: Authentication headers

    Returns:
        Inventory summary data
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{backend_url}/api/agent/operate",
            headers=headers,
            json={
                "operation": "generate_report",
                "resource": "inventory",
                "parameters": {"report_type": "inventory_summary"},
                "reasoning": "Generate inventory summary report",
            },
        )
        response.raise_for_status()
        return response.json()


async def generate_material_plan_summary(
    backend_url: str,
    headers: dict,
) -> dict[str, Any]:
    """Generate material plan summary report

    Args:
        backend_url: Backend API base URL
        headers: Authentication headers

    Returns:
        Material plan summary data
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{backend_url}/api/agent/operate",
            headers=headers,
            json={
                "operation": "generate_report",
                "resource": "material_plan",
                "parameters": {"report_type": "material_plan_summary"},
                "reasoning": "Generate material plan summary report",
            },
        )
        response.raise_for_status()
        return response.json()

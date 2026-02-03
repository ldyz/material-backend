"""
Inventory Tools for MCP Server
Provides tools for querying and analyzing inventory data
"""

from typing import Any
import httpx


async def query_materials(
    backend_url: str,
    headers: dict,
    search: str = "",
    limit: int = 10,
) -> dict[str, Any]:
    """Query material master data

    Args:
        backend_url: Backend API base URL
        headers: Authentication headers
        search: Search keyword for material name or specification
        limit: Maximum number of results

    Returns:
        Query results as dictionary
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{backend_url}/api/agent/query",
            headers=headers,
            json={
                "question": f"Search materials: {search}",
                "limit": limit,
            },
        )
        response.raise_for_status()
        return response.json()


async def analyze_inventory(
    backend_url: str,
    headers: dict,
    question: str,
) -> dict[str, Any]:
    """Analyze inventory data

    Args:
        backend_url: Backend API base URL
        headers: Authentication headers
        question: Analysis question

    Returns:
        Analysis results
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{backend_url}/api/agent/operate",
            headers=headers,
            json={
                "operation": "analyze",
                "resource": "inventory",
                "parameters": {"question": question},
                "reasoning": f"Inventory analysis: {question}",
            },
        )
        response.raise_for_status()
        return response.json()


async def get_stock_alerts(
    backend_url: str,
    headers: dict,
) -> dict[str, Any]:
    """Get low stock alerts

    Args:
        backend_url: Backend API base URL
        headers: Authentication headers

    Returns:
        Stock alerts data
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{backend_url}/api/agent/operate",
            headers=headers,
            json={
                "operation": "query",
                "resource": "stock",
                "parameters": {"low_stock_alert": True},
                "reasoning": "Get stock alerts",
            },
        )
        response.raise_for_status()
        return response.json()

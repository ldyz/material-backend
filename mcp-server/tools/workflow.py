"""
Workflow Tools for MCP Server
Provides tools for workflow operations and approvals
"""

from typing import Any
import httpx


async def get_pending_tasks(
    backend_url: str,
    headers: dict,
    limit: int = 20,
) -> dict[str, Any]:
    """Get pending workflow tasks

    Args:
        backend_url: Backend API base URL
        headers: Authentication headers
        limit: Maximum number of results

    Returns:
        Pending tasks data
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{backend_url}/api/agent/operate",
            headers=headers,
            json={
                "operation": "query",
                "resource": "workflow",
                "parameters": {"status": "pending", "limit": limit},
                "reasoning": "Get pending workflow tasks",
            },
        )
        response.raise_for_status()
        return response.json()


async def approve_workflow(
    backend_url: str,
    headers: dict,
    task_id: int,
    remark: str = "",
) -> dict[str, Any]:
    """Approve a workflow task

    Args:
        backend_url: Backend API base URL
        headers: Authentication headers
        task_id: Workflow task ID
        remark: Approval remarks

    Returns:
        Approval result
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{backend_url}/api/agent/workflow",
            headers=headers,
            json={
                "task_id": task_id,
                "action": "approve",
                "remark": remark,
            },
        )
        response.raise_for_status()
        return response.json()


async def reject_workflow(
    backend_url: str,
    headers: dict,
    task_id: int,
    reason: str = "",
) -> dict[str, Any]:
    """Reject a workflow task

    Args:
        backend_url: Backend API base URL
        headers: Authentication headers
        task_id: Workflow task ID
        reason: Rejection reason

    Returns:
        Rejection result
    """
    async with httpx.AsyncClient() as client:
        response = await client.post(
            f"{backend_url}/api/agent/workflow",
            headers=headers,
            json={
                "task_id": task_id,
                "action": "reject",
                "remark": reason,
            },
        )
        response.raise_for_status()
        return response.json()


async def get_workflow_history(
    backend_url: str,
    headers: dict,
    instance_id: int,
) -> dict[str, Any]:
    """Get workflow approval history

    Args:
        backend_url: Backend API base URL
        headers: Authentication headers
        instance_id: Workflow instance ID

    Returns:
        Workflow history
    """
    async with httpx.AsyncClient() as client:
        response = await client.get(
            f"{backend_url}/api/workflow-instances/{instance_id}/approvals",
            headers=headers,
        )
        response.raise_for_status()
        return response.json()

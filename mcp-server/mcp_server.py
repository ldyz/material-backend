#!/usr/bin/env python3
"""
MCP Server for Material Management System
Connects Claude Desktop to the Material Management Backend
"""

import argparse
import asyncio
import json
import os
import sys
from typing import Any, Sequence

import httpx
import yaml
from dotenv import load_dotenv
from mcp.server import Server
from mcp.server.models import InitializationOptions
from mcp.types import (
    CallToolRequest,
    CallToolResult,
    ListToolsRequest,
    ListToolsResult,
    Tool,
    TextContent,
)

# Load environment variables
load_dotenv()

# Load configuration
CONFIG_PATH = os.path.join(os.path.dirname(__file__), "config.yaml")

# Global configuration
config = {}
client = None

# Initialize MCP Server
app = Server("material-management-mcp")


def load_config(config_path: str = CONFIG_PATH) -> dict:
    """Load configuration from YAML file"""
    global config
    with open(config_path, "r") as f:
        config = yaml.safe_load(f)
    return config


def get_auth_headers() -> dict:
    """Get authentication headers for backend requests"""
    token = os.getenv("JWT_TOKEN")
    if not token:
        # Try to get token from backend
        backend_url = config.get("backend", {}).get("base_url", "http://localhost:8088")
        username = os.getenv("BACKEND_USERNAME", "admin")
        password = os.getenv("BACKEND_PASSWORD", "admin")

        # In production, you should get this from a secure source
        # For now, we'll use a placeholder
        token = os.getenv("JWT_SECRET", "your-jwt-token")

    return {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json",
        "X-Agent-ID": "claude-desktop-mcp-v1",
    }


async def backend_request(method: str, endpoint: str, data: dict = None) -> dict:
    """Make a request to the backend API"""
    backend_url = config.get("backend", {}).get("base_url", "http://localhost:8088")
    url = f"{backend_url}/api/agent{endpoint}"

    async with httpx.AsyncClient(timeout=30.0) as client:
        if method == "GET":
            response = await client.get(url, headers=get_auth_headers(), params=data)
        elif method == "POST":
            response = await client.post(url, headers=get_auth_headers(), json=data)
        else:
            raise ValueError(f"Unsupported method: {method}")

        response.raise_for_status()
        return response.json()


@app.list_tools()
async def list_tools() -> list[Tool]:
    """List available MCP tools"""
    tools_config = config.get("tools", [])
    tools = []

    for tool_config in tools_config:
        if not tool_config.get("enabled", True):
            continue

        tool = Tool(
            name=tool_config["name"],
            description=tool_config["description"],
            inputSchema={
                "type": "object",
                "properties": {},
            },
        )

        # Add parameters from config
        if "parameters" in tool_config:
            for param in tool_config["parameters"]:
                tool.inputSchema["properties"][param["name"]] = {
                    "type": param.get("type", "string"),
                    "description": param.get("description", ""),
                }
                if param.get("required", False):
                    tool.inputSchema.setdefault("required", []).append(param["name"])

        tools.append(tool)

    return tools


@app.call_tool()
async def call_tool(name: str, arguments: dict) -> list[TextContent]:
    """Handle tool calls"""
    try:
        if name == "query_materials":
            return await query_materials(arguments)
        elif name == "analyze_inventory":
            return await analyze_inventory(arguments)
        elif name == "create_material_plan":
            return await create_material_plan(arguments)
        elif name == "get_stock_alerts":
            return await get_stock_alerts(arguments)
        elif name == "approve_workflow":
            return await approve_workflow(arguments)
        elif name == "get_pending_tasks":
            return await get_pending_tasks(arguments)
        else:
            return [TextContent(type="text", text=f"Unknown tool: {name}")]
    except Exception as e:
        return [TextContent(type="text", text=f"Error: {str(e)}")]


async def query_materials(args: dict) -> list[TextContent]:
    """Query material master data"""
    search = args.get("search", "")
    limit = args.get("limit", 10)

    result = await backend_request(
        "POST",
        "/query",
        {
            "question": f"Search materials: {search}",
            "limit": limit,
        },
    )

    return [
        TextContent(
            type="text",
            text=json.dumps(result, indent=2, ensure_ascii=False),
        )
    ]


async def analyze_inventory(args: dict) -> list[TextContent]:
    """Analyze inventory data"""
    question = args.get("question", "")

    result = await backend_request(
        "POST",
        "/operate",
        {
            "operation": "analyze",
            "resource": "inventory",
            "parameters": {
                "question": question,
            },
            "reasoning": f"User requested inventory analysis: {question}",
        },
    )

    return [
        TextContent(
            type="text",
            text=json.dumps(result, indent=2, ensure_ascii=False),
        )
    ]


async def create_material_plan(args: dict) -> list[TextContent]:
    """Create a material plan"""
    project_id = args.get("project_id")
    items = args.get("items", [])
    remark = args.get("remark", "")

    result = await backend_request(
        "POST",
        "/operate",
        {
            "operation": "create_material_plan",
            "resource": "material_plan",
            "parameters": {
                "project_id": project_id,
                "items": items,
                "remark": remark,
            },
            "reasoning": f"AI generated material plan for project {project_id}",
        },
    )

    return [
        TextContent(
            type="text",
            text=json.dumps(result, indent=2, ensure_ascii=False),
        )
    ]


async def get_stock_alerts(args: dict) -> list[TextContent]:
    """Get low stock alerts"""
    result = await backend_request(
        "POST",
        "/operate",
        {
            "operation": "query",
            "resource": "stock",
            "parameters": {
                "low_stock_alert": True,
            },
            "reasoning": "Get stock alerts for AI monitoring",
        },
    )

    return [
        TextContent(
            type="text",
            text=json.dumps(result, indent=2, ensure_ascii=False),
        )
    ]


async def approve_workflow(args: dict) -> list[TextContent]:
    """Approve a workflow task"""
    task_id = args.get("task_id")
    remark = args.get("remark", "")

    result = await backend_request(
        "POST",
        "/workflow",
        {
            "task_id": task_id,
            "action": "approve",
            "remark": remark,
        },
    )

    return [
        TextContent(
            type="text",
            text=json.dumps(result, indent=2, ensure_ascii=False),
        )
    ]


async def get_pending_tasks(args: dict) -> list[TextContent]:
    """Get pending workflow tasks"""
    result = await backend_request(
        "POST",
        "/operate",
        {
            "operation": "query",
            "resource": "workflow",
            "parameters": {
                "status": "pending",
            },
            "reasoning": "Get pending tasks for AI review",
        },
    )

    return [
        TextContent(
            type="text",
            text=json.dumps(result, indent=2, ensure_ascii=False),
        )
    ]


async def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(description="MCP Server for Material Management")
    parser.add_argument(
        "--config",
        default=CONFIG_PATH,
        help="Path to configuration file",
    )
    args = parser.parse_args()

    # Load configuration
    load_config(args.config)

    # Run MCP server
    from mcp.server.stdio import stdio_server

    async with stdio_server() as (read_stream, write_stream):
        await app.run(
            read_stream,
            write_stream,
            InitializationOptions(
                server_name=config.get("server", {}).get("name", "material-management-mcp"),
                server_version=config.get("server", {}).get("version", "1.0.0"),
                capabilities=app.get_capabilities(
                    notification_options=None,
                    experimental_capabilities={},
                ),
            ),
        )


if __name__ == "__main__":
    asyncio.run(main())

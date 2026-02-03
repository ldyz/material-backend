"""
Prompt Templates for MCP Server
Provides reusable prompt templates for AI Agent interactions
"""

from typing import Any


# Query Templates
QUERY_MATERIALS_TEMPLATE = """
User is asking about materials. Please use the query_materials tool with these parameters:
- Search term: {search}
- Limit: {limit}

User question: {question}
"""

ANALYZE_INVENTORY_TEMPLATE = """
User is requesting inventory analysis. Please use the analyze_inventory tool.

Analysis question: {question}

Context: {context}
"""

# Action Templates
CREATE_MATERIAL_PLAN_TEMPLATE = """
User wants to create a material plan. Please use the create_material_plan tool.

Project ID: {project_id}
Items: {items}
Reasoning: {reasoning}

Additional notes: {remark}
"""

APPROVE_WORKFLOW_TEMPLATE = """
User wants to approve a workflow task. Please use the approve_workflow tool.

Task ID: {task_id}
Reasoning: {reasoning}

Additional remarks: {remark}
"""

# System Prompts
SYSTEM_INVENTORY_ASSISTANT = """
You are an Inventory Management Assistant for a Material Management System.
You can help users:
- Query material information
- Analyze inventory levels
- Get stock alerts
- Generate inventory reports

Always explain your reasoning before taking actions.
"""

SYSTEM_WORKFLOW_ASSISTANT = """
You are a Workflow Management Assistant for a Material Management System.
You can help users:
- View pending workflow tasks
- Approve or reject workflow tasks
- Check workflow status
- View workflow history

Always explain your reasoning before taking actions.
For sensitive operations like approvals, clearly state what you're doing and why.
"""

SYSTEM_MATERIAL_PLANNING_ASSISTANT = """
You are a Material Planning Assistant for a Material Management System.
You can help users:
- Create material plans
- Update existing plans
- Analyze plan requirements
- Generate planning reports

Always explain your reasoning before taking actions.
Ensure all required information is collected before creating plans.
"""


def get_template(template_name: str, **kwargs: Any) -> str:
    """Get a formatted prompt template

    Args:
        template_name: Name of the template
        **kwargs: Variables to format into the template

    Returns:
        Formatted template string
    """
    templates = {
        "query_materials": QUERY_MATERIALS_TEMPLATE,
        "analyze_inventory": ANALYZE_INVENTORY_TEMPLATE,
        "create_material_plan": CREATE_MATERIAL_PLAN_TEMPLATE,
        "approve_workflow": APPROVE_WORKFLOW_TEMPLATE,
        "system_inventory": SYSTEM_INVENTORY_ASSISTANT,
        "system_workflow": SYSTEM_WORKFLOW_ASSISTANT,
        "system_planning": SYSTEM_MATERIAL_PLANNING_ASSISTANT,
    }

    template = templates.get(template_name)
    if template is None:
        raise ValueError(f"Unknown template: {template_name}")

    return template.format(**kwargs)

"""
MkDocs hook to fix line endings for proper rendering.
Automatically adds two spaces at the end of lines that need line breaks.
"""
import re


def on_page_markdown(markdown, page, config, **kwargs):
    """
    Fix line endings in markdown content for proper MkDocs rendering.
    Adds two spaces at the end of lines that need line breaks.
    """
    lines = markdown.split('\n')
    processed_lines = []
    in_code_block = False
    
    for i, line in enumerate(lines):
        stripped = line.strip()
        
        # Track code blocks
        if stripped.startswith('```'):
            in_code_block = not in_code_block
            processed_lines.append(line)
            continue
            
        # Skip processing inside code blocks

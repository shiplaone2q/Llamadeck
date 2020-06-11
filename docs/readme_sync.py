"""
MkDocs hook to sync content from README.md to docs/index.md
"""
import re
import os


def on_page_markdown(markdown, page, config, **kwargs):
    """Process markdown content before rendering"""
    # Only process the index.md file
    if page.file.src_path != 'index.md':
        return markdown

    # Get the path to README.md (relative to mkdocs.yml)
    readme_path = os.path.join(os.path.dirname(config['config_file_path']), 'README.md')

    if not os.path.exists(readme_path):
        print(f"Warning: README.md not found at {readme_path}")
        return markdown

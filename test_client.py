#!/usr/bin/env python3
"""
Simple Python script to interact with local LLM server's OpenAI-compatible API
"""

import requests
import sys

# Local LLM server configuration
BASE_URL = "http://localhost:8080"
API_KEY = None
MODEL_NAME = None

def get_models():
    """Fetch available models from /v1/models endpoint"""
    headers = {}
    if API_KEY:
        headers["Authorization"] = f"Bearer {API_KEY}"

    try:
        response = requests.get(f"{BASE_URL}/v1/models", headers=headers, timeout=10)
        response.raise_for_status()
        return response.json()["data"]
    except Exception as e:
        print(f"Error fetching models: {e}")
        return []

def send_message(message):
    """
    Send a message to local LLM server API

    Args:
        message (str): The message to send

    Returns:
        str: The AI response or error message
    """

    headers = {
        "Content-Type": "application/json",
    }

    if API_KEY:
        headers["Authorization"] = f"Bearer {API_KEY}"

    data = {
        "model": MODEL_NAME,
        "messages": [
            {
                "role": "user",
                "content": message
            }
        ],
        "temperature": 0.7,
        "max_tokens": 1000,
        "stream": False,
    }
    
    response = requests.post(f"{BASE_URL}/v1/chat/completions", headers=headers, json=data, timeout=60)
    response.raise_for_status()
    return response.json()["choices"][0]["message"]["content"]

def interactive_mode():
    """Run in interactive mode for continuous conversation"""
    global BASE_URL, API_KEY, MODEL_NAME

    # Get base URL
    url_input = input(f"Base URL [{BASE_URL}]: ").strip()
    if url_input:
        BASE_URL = url_input

    # Get API key (optional)
    key_input = input("API key (optional): ").strip()
    if key_input:
        API_KEY = key_input

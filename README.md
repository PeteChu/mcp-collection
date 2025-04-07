# Repository Overview

This repository is a collection of Model Context Protocol (MCP) implemented in both Rust and Go. Each MCP exposes its functionality through an MCP/RMCP interface, allowing tools to be registered and executed remotely. The repository currently contains the following projects:

- **Calculator Service (Rust):**  
  A simple calculator that implements basic arithmetic operations such as addition and subtraction.

- **Weather Service (Rust):**  
  A weather service that fetches geographical coordinates and weather information from the OpenWeatherMap API. It includes functions to retrieve the current weather and a 5-day forecast.

- **Metal Price Service (Go):**  
  A service that interacts with a metal price API. It supports functionalities including:
  - Getting the current date.
  - Listing supported currency symbols.
  - Retrieving live rates, historical rates, and OHLC (Open, High, Low, Close) data.
  - Querying exchange rates over a specific timeframe.

# Introduction

## Purpose of This Document

This document serves as the introduction to the architectural documentation of the **Esp Monitor** platform.

Its purpose is to introduce the reader to the project's goals, scope, and overall concept. Detailed descriptions of the architecture, design principles, and individual components are provided in the subsequent chapters.

---

## What Is Esp Monitor?

**Esp Monitor** is a self-hosted platform for the initial provisioning, monitoring, and remote management of devices based on ESP-family microcontrollers.

The project is intended for equipment installers, system integrators, embedded software developers, and administrators who need to manage large numbers of devices without relying on proprietary cloud services.

The primary objective of the platform is to provide a unified solution for managing the entire device lifecycle—from the initial connection to a local network through day-to-day operation, firmware updates, and diagnostics.

---

## Why the Platform Exists

Many existing IoT management solutions are tightly coupled to the manufacturer's cloud infrastructure or depend on proprietary communication protocols.

Esp Monitor follows a different approach. The platform is designed to be fully deployable within an organization's own infrastructure—or even on a single local computer—without requiring connectivity to external services.

This approach allows organizations to retain complete control over their infrastructure, device management processes, and the long-term evolution of the system.

---

## Key Characteristics

The platform is designed with particular emphasis on the following principles:

- Self-hosted deployment;
- Open technologies and standard networking protocols;
- Minimal external dependencies;
- Long-term maintainability and extensibility.

These characteristics define the overall direction of the project. The engineering principles behind them are discussed in detail in the following chapter.

---

## Intended Audience

The platform is primarily intended for:

- equipment installers;
- ESP-based device developers;
- system integrators;
- small companies managing distributed device fleets;
- educational institutions;
- makers and IoT enthusiasts.

Although the platform is suitable for commercial deployments, its architecture remains simple enough for home laboratories, research projects, and educational environments.

---

## Documentation Structure

The architectural documentation is organized into independent chapters, each covering a specific architectural responsibility within the platform.

Throughout the documentation, a consistent architectural classification is used, including the following areas:

- Platform Core
- External Interfaces
- Persistence
- Configuration
- Workers
- Infrastructure Services
- Operational Infrastructure
- Reference Implementations

This structure enables readers to explore the architecture progressively while making the documentation easier to maintain as the platform evolves.

---

## Next Steps

After completing this introduction, readers are encouraged to continue with the next document:

**02 — Design Philosophy**

This chapter describes the engineering principles that drive the architectural decisions and guide the long-term evolution of the platform.
# Article Management System

A full-stack article management system built using Go for the backend REST API and Next.js for the frontend dashboard.

This project demonstrates a separated architecture between backend services and frontend interface, focusing on API design, data handling, and dashboard interaction.

---

## 📁 Project Structure

This repository is a monorepo containing two main services:

```
/article-services-go       → Backend (Go REST API)
/article-dashboard-nextjs → Frontend (Next.js Dashboard)
```

---

## 🛠️ Tech Stack

### Backend

* Go (Golang)
* REST API
* Go Modules

### Frontend

* Next.js
* React
* JavaScript / TypeScript
* Axios

---

## 🚀 Getting Started

This project runs in two separate services. You need to run both backend and frontend in different terminals.

---

## 1. Backend (Go Service)

### Navigate to backend folder

```bash
cd article-services-go
```

### Install dependencies

```bash
go mod tidy
```

### Run backend server

```bash
go run .
```

Backend will run at:

```
http://localhost:8080
```

---

## 2. Frontend (Next.js Dashboard)

### Open a new terminal

### Navigate to frontend folder

```bash
cd article-dashboard-nextjs
```

### Install dependencies

```bash
npm install
```

### Run development server

```bash
npm run dev
```

Frontend will run at:

```
http://localhost:3000
```

---

## ✨ Features

* Article CRUD (Create, Read, Update, Delete)
* RESTful API built with Go
* Dashboard interface using Next.js
* Separation of backend and frontend services
* Simple and scalable architecture

---

## 🎯 Purpose of Project

This project was built to practice:

* Building REST APIs using Go
* Consuming APIs using Next.js frontend
* Full-stack application architecture
* Monorepo project structure
* Basic content management system design

---

## 📌 Notes

This is a learning project focused on full-stack integration between Go and Next.js. It is not intended for production use.

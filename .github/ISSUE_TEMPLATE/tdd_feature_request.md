---
name: 🧪 TDD Feature Request
about: 提出一個新功能或任務，並定義對應的測試策略 (Strict TDD)
title: '[TDD] <請簡述功能名稱>'
labels: ['status: planning']
assignees: ''
---

## 🎯 Objective (目標)
> 本任務旨在完成 [Backend/Agent] 的 _____________ 功能。

## 🧪 TDD Strategy (核心測試策略)
**1. Test Scenario (測試場景):**
**2. Test Type (測試類型):**
- [ ] **Unit Test (Mocking)**: 用於 Agent Runner 邏輯 (不依賴真實 HPL) 
- [ ] **Integration Test (Testcontainers)**: 用於 Core DB 存取與 Ranking 計算 
- [ ] **Golden File Test**: 用於 HPL Log Parser 解析驗證 
- [ ] **API Contract Test**: 用於驗證 HTTP Request/Response 格式 

**3. Expected Behavior (預期行為/驗收標準):**
- **Input / Setup:** (例如: `Input Log: "HPL result: NaN"`)
- **Expected Output:** (例如: `Parser should throw InvalidLogFormatException`)

## 🛠 Implementation Plan (實作計畫)
- [ ] 定義介面 (Interface/DTO)
- [ ] 撰寫測試代碼 (The "Red" Phase)
- [ ] 實作最小功能代碼 (The "Green" Phase)
- [ ] 重構 (Refactor)

## 👤 Owner & Role
- **Component:** [Core Backend / Agent Client]
- **Assignee:**
  - [ ] 109704065 (Backend Lead: API & DB)
  - [ ] 113550064 (Agent Lead: Parser & Runner)

---
*Remember: No implementation may be merged without a failing test first.* 
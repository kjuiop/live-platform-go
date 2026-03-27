---
name: question
description: 질문에 대한 답변을 docs/claude/question/ 아래에 MD 파일로 저장합니다
---

`$ARGUMENTS`에 담긴 질문에 답변하고, 그 내용을 `docs/claude/question/` 아래에 MD 파일로 저장합니다.

## 1. 파일명 결정

- 질문 내용을 요약해 짧은 영문 슬러그를 만듭니다
- 형식: `docs/claude/question/{슬러그}.md`
- 예시: `docs/claude/question/adapter-vs-controller.md`
- `docs/claude/question/` 디렉터리가 없으면 생성합니다

## 2. 답변 작성

아래 구조로 문서를 작성합니다:

```markdown
# Q. {질문 원문}

> 작성일: {오늘 날짜}

## 답변

{질문에 대한 명확하고 구체적인 답변}

---

## 관련 개념

{필요한 경우 추가 배경 지식이나 관련 개념 설명}

---

## 참고 링크

- [{제목}]({URL})
```

작성 규칙:
- 답변은 이 프로젝트(Go 1.25, Gin, Gorilla WebSocket, Kafka, Redis, Clean Architecture) 맥락에 맞게 작성합니다
- 코드 예시가 필요하면 구체적인 변수명과 패키지 경로를 포함합니다
- "관련 개념"과 "참고 링크"는 내용이 없으면 생략합니다

## 3. 결과 출력

- 생성된 문서 경로를 알려줍니다
- 답변 핵심 요약을 2~3줄로 출력합니다

주의사항:
- Go 코드를 직접 프로젝트 파일에 작성하거나 수정하지 않습니다 — 문서만 작성합니다
- `$ARGUMENTS`가 비어있으면 사용자에게 질문 내용을 물어봅니다

---
name: pr-create
description: 현재 브랜치의 커밋을 분석해 PR을 생성합니다
---

현재 브랜치의 PR을 생성합니다.

다음 순서로 진행하세요:

1. 아래 명령어를 병렬로 실행해 현재 상태를 파악하세요:
   - `git status`
   - `git log main..HEAD --oneline`
   - `git diff main...HEAD --stat`

2. 커밋 내역과 변경 파일을 분석해 PR 제목과 본문을 작성하세요.
   - 제목: 커밋 메시지를 그대로 쓰지 말고, 변경 내용을 파악해 **직접 의미 있는 제목**을 작성합니다 (70자 이내)
   - 본문: 반드시 `.github/pull_request_template.md` 파일을 읽어 **HTML 주석(`<!-- ... -->`)을 포함한 전체 내용을 그대로** 사용합니다
     - 템플릿의 모든 섹션(Summary, Issue, Why, What, How, Test)을 유지합니다
     - 상단 Copilot Review Instruction 주석 블록도 반드시 포함합니다
     - 각 섹션을 커밋 내역과 변경 파일 분석 결과로 채웁니다

3. 원격 브랜치 존재 여부를 확인하세요:
   - 원격 브랜치가 없으면 **즉시 중단**하고 사용자에게 알립니다 (push 금지)
   - 원격 브랜치가 있으면 PR을 생성하세요:
   ```
   gh pr create --title "..." --body "$(cat <<'EOF'
   <!-- 템플릿 전체 내용 (주석 포함) -->
   EOF
   )"
   ```

4. 생성된 PR URL을 출력하세요.

주의사항:
- force push 금지
- push 금지 — 커밋과 원격 브랜치 생성은 사용자가 직접 한다
- main 브랜치에서 직접 실행 시 경고 후 중단
- 변경사항이 없으면 PR 생성하지 않음
- 원격 브랜치가 없으면 PR 생성하지 않고 중단

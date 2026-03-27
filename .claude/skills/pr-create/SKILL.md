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
   - 제목: 70자 이내
   - 본문: 프로젝트의 `.github/pull_request_template.md` 형식을 그대로 따릅니다

3. 원격 브랜치가 없으면 push 후 PR을 생성하세요:
   ```
   gh pr create --title "..." --body "..."
   ```

4. 생성된 PR URL을 출력하세요.

주의사항:
- force push 금지
- main 브랜치에서 직접 실행 시 경고 후 중단
- 변경사항이 없으면 PR 생성하지 않음

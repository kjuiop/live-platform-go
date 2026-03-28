---
name: pr-review
description: 현재 브랜치의 열린 PR을 Go 서버 관점에서 리뷰하고 코멘트를 등록합니다
---

현재 브랜치의 열린 PR을 찾아 코드 리뷰를 작성합니다.

다음 순서로 진행하세요:

1. 현재 브랜치의 PR 번호를 확인합니다:
   ```
   gh pr view --json number,title,body,url
   ```

2. PR의 변경 diff를 가져옵니다:
   ```
   gh pr diff
   ```

3. 변경된 파일 목록을 확인합니다:
   ```
   gh pr view --json files
   ```

4. 아래 관점에서 리뷰를 수행합니다:
   - **동시성/경쟁 조건**: Goroutine, Channel, sync 패키지 사용의 안전성
   - **Kafka 처리**: Producer/Consumer 에러 처리, 오프셋 관리, 재시도 로직
   - **WebSocket**: 연결 종료 처리, 읽기/쓰기 타임아웃, 동시 쓰기 방지
   - **에러 처리**: 에러 무시 여부, 적절한 래핑 (`fmt.Errorf` + `%w`)
   - **리소스 누수**: defer close, context 취소 전파
   - **성능**: 불필요한 alloc, blocking 호출 위치

5. PR에 이미 달린 **기존 코멘트(Copilot 포함)**를 가져와 중복을 방지합니다:
   ```
   gh api repos/{owner}/{repo}/pulls/{pr_number}/comments --jq '[.[] | {path: .path, line: .line, body: .body}]'
   gh api repos/{owner}/{repo}/issues/{pr_number}/comments --jq '[.[] | {body: .body}]'
   ```
   - 기존 코멘트와 **동일하거나 유사한 지적**(같은 파일·라인, 같은 문제)은 등록하지 않습니다
   - 이미 Copilot이 지적한 내용은 건너뜁니다

6. 리뷰 코멘트를 **라인별로** PR에 등록합니다:
   - 전체 요약 댓글 하나로 올리지 않고, 지적 사항마다 해당 파일·라인에 직접 코멘트를 답니다
   - PR의 최신 커밋 SHA를 먼저 확인합니다:
     ```
     gh pr view --json commits --jq '.commits[-1].oid'
     ```
   - 라인 코멘트는 GitHub REST API로 등록합니다:
     ```
     gh api repos/{owner}/{repo}/pulls/{pr_number}/comments \
       --method POST \
       --field body="🤖 Claude Review

코멘트 내용" \
       --field commit_id="<커밋 SHA>" \
       --field path="파일 경로" \
       --field line=<라인 번호> \
       --field side="RIGHT"
     ```
   - `owner`와 `repo`는 `gh repo view --json owner,name`으로 확인합니다
   - `line`은 diff 기준 **변경된 라인(+줄)**의 번호를 사용합니다
   - 지적할 라인이 없는 종합 의견은 PR review comment로 남깁니다:
     ```
     gh pr review --comment --body "🤖 Claude Review

..."
     ```

7. 리뷰 결과 요약을 출력합니다 (LGTM / 수정 필요 항목 / 중복으로 건너뛴 항목).

주의사항:
- 열린 PR이 없으면 안내 메시지 출력 후 중단
- Go 코드를 직접 수정하지 않음 — 리뷰 코멘트만 작성
- 댓글 하나로 묶어서 올리지 않음 — 반드시 지적 사항별로 해당 라인에 개별 코멘트 등록

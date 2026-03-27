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

5. 리뷰 코멘트를 PR에 등록합니다:
   - 전체 리뷰 요약: `gh pr review --comment --body "..."`
   - 라인별 코멘트가 필요하면 GitHub CLI 또는 API를 사용합니다

6. 리뷰 결과 요약을 출력합니다 (LGTM / 수정 필요 항목).

주의사항:
- 열린 PR이 없으면 안내 메시지 출력 후 중단
- Go 코드를 직접 수정하지 않음 — 리뷰 코멘트만 작성

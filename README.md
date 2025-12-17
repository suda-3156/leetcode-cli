# leetcode-cli

A command-line interface (CLI) tool for generating LeetCode problem files.

This CLI is not intended for submitting solutions to LeetCode or searching for problems. Instead, it is designed solely for generating problem files with templates.
This is because you can easily submit solutions and search for problems on the LeetCode website itself, where you can also find more information.

Currently, only Python3 and Go are supported for extracting type definitions.

## TODO

- [ ] Support more languages (Currently, only Python3 and Go are supported for extracting type definitions)
- [ ] Support custom templates for each language

## Usage

```bash
lcli generate <keyword>
```

- `<keyword>`: The keyword to search for the problem (e.g., problem title or ID).

## Installation

Requires Go 1.25.4 or higher:

```bash
go install github.com/suda-3156/leetcode-cli@latest # With binary name `leetcode-cli`
```

## Example

When you execute the following command:

```sh
lcli 19 # And choose python3 as language
```

The following file will be generated to `./yyyy-mm-dd/19.remove-nth-node-from-end-of-list.py`:

```py
# 2025-12-13
# 19. Remove Nth Node From End of List
# Python3


from typing import TypedDict, Any, List


class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next

# Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class Solution:
    def removeNthFromEnd(self, head: Optional[ListNode], n: int) -> Optional[ListNode]:

        pass


TestCase = TypedDict("TestCase", {"input": Any, "want": Any})


def main():
    test_cases: List[TestCase] = [
        # Add your test cases here
        {"input": {}, "want": ""}
    ]

    s = Solution()

    for tc in test_cases:
        result = s.removeNthFromEnd(tc["input"])
        if result != tc["want"]:
            print(f"Test failed for input {tc["input"]}: got {result}, want {tc["want"]}")
        else:
            print(f"Test passed for input {tc["input"]}")


if __name__ == "__main__":
    main()
```

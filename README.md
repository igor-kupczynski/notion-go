# notion-go

An experiment in creating a go client for the notion
api, [currently in beta](https://twitter.com/NotionHQ/status/1392883313757409284).

The goal is to use a notion database to drive some content on [my blog](https://kupczynski.info)

## Experiment status

**2021-03-22** 👎 I can't use the notion beta API to run content of my blog. Notion pages consist of blocks, and in the
beta API [only text-like blocks are currently supported](https://developers.notion.com/reference/block)

> At present the API only supports text-like block types which are listed in the reference below. All other block types will continue to appear in the structure, but only contain a type set to "unsupported".

Specifically, this means that the code blocks are `"unsupported"`, and as you can imagine there's a lot of code blocks
on a technical blog.

I'll be happy to revisit once the API is extended.

Currently the best bet seem to be the [unofficial API go client](https://github.com/kjk/notionapi).

## Implementation status

* Databases
    - [x] Retrieve a database
    - [x] Query a database
    - [x] List databases
    - ⚠️ not all properties and filter types are implemented

* Pages
    - [ ] Retrieve a page
    - [ ] Create a page
    - [ ] Update page properties

* Blocks
    - [ ] Retrieve block children
    - [ ] Append block children

* Users
    - [ ] Retrieve a user
    - [ ] List all users

* Search
    - [ ] Search


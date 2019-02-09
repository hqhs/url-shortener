Classic service example: URL shortener
===

## Features

- simple analytics

## TODO: Code coverage

## Database choice

Let's say average long URL size is 2kb. Additionally we want to record
statistics, such as Unix timestamp and IP request address for each click. Keep
in mind that we may add more statistics in future. 10 bytes for timestamp and 16
bytes for IP (with IPv6 support) sounds reasonable. Some might argue what
statistics is not required, but we're living in word of data driven advertising
and loosing opportunity to collect data in my opinion is not realistic.

bit.ly statistics is ([link]( http://highscalability.com/blog/2014/7/14/bitly-lessons-learned-building-a-distributed-system-that-han.html )): 
  - 6 * 10^9 clicks/month
  - 6 * 10^8 shortens/month,
  - avarage click/link is 10, but max(click per link) could be 10^6 or even more.
 
 6 * 10^9 * 26 bytes + 6 * 10^8 * 2 kb = 1.35 terabytes/month
 
 Obviously we need to choose database wisely. Out data structure is fixed,
 with small quantities, we don't need atomicity.
 
|                              | postgres/mysql | redis cluster | mongo   | hadoop  | clickhouse |
|:-----------------------------|----------------|:--------------|:--------|:--------|:-----------|
| structured                   | **yes**        | no            | **yes** | **yes** | **yes**    |
| easy to scale                | no             | **yes**       | **yes** | **yes** | **yes**    |
| easy to get up and running   | **yes**        | **yes**       | **yes** | no      | **yes**    |
| fast with our data structure | no             | **yes**       | no      | **yes** | **yes**    |
| easy to maintain             | **yes**        | **yes**       | no      | no      | **yes**    |
| reliable                     | **yes**        | no            | no      | **yes** | **yes**    |

 **NOTE**: this is simple and very opinionated binary analysis based on my
 experience and limited knowledge. Purpose of this comparison is consider
 different options with given conditions. Your choice may differ, but feel free
 to argue!

 Remember, what statistics above taken from startup which is **primary** product
 is URL shortening & analytics service. Just copy-pasting their architecture
 (and choosing hadoop) solutions doesn't fit because we shouldn't **expect**
 such traffic. Our goal is to find balance between scale fees and time to get
 service up and running fast.
 
## URL shortener algorithm choice

## statistics: codebase size & time spent

## possible improvements


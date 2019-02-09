Classic service example: URL shortener
===

### Features

- simple analytics

### TODO: Code coverage

### Database choice

Let's say avarage long URL size is 2kb, using base64 encoding and starting with
6 chars for short url, we get 64^6 urls, which is aprox 10^11 records.
Additionally we want to record statistics, such as unix timestamp and ip request
address for each click. Keep in mind that we may add more statistics in future. 
10 bytes for timestamp and 16 bytes for ip (for ipv6 support) sounds reasonable.
Some might argue what statistics is not required, but we're living in word of
data driven advertising and loosing opportunity to collect data in my opinion is
not realistic.

bit.ly stats is ([ link ]( http://highscalability.com/blog/2014/7/14/bitly-lessons-learned-building-a-distributed-system-that-han.html )): 
 6 * 10^9 clicks/month, 6 * 10^8 shortens/month, avg click/link is 10, but
 max(click per link) could be 10^6 or more.
 
 6 * 10^9 * 26 bytes + 6 * 10^8 * 2 kb = 1.35 terabytes/month
 
 Obviosly we need to choose database wisely. Out data structure is fixed,
 with small quantities, we dont need atomicity.
 
|                              | postgres/mysql | redis cluster | mongo   | hadoop  | clickhouse |
|:-----------------------------|----------------|:--------------|:--------|:--------|:-----------|
| structured                   | **yes**        | no            | **yes** | **yes** | **yes**    |
| easy to scale                | no             | **yes**       | **yes** | **yes** | **yes**    |
| easy to get up and running   | **yes**        | **yes**       | **yes** | no      | **yes**    |
| fast with out data structure | no             | **yes**       | no      | **yes** | **yes**    |
| easy to maintain             | **yes**        | **yes**       | no      | no      | **yes**    |
| reliable                     | **yes**        | no            | no      | **yes** | **yes**    |

 NOTE: this is simple and very opinionated binary analysis based on my
 expirience and limited knowledge. Purpose of this comparison is consider
 different options with given conditions. Your choice may differ, but feel free
 to argue!

 Remember, what statistics above taken from startup which is **primary** product
 is url shortening & analytics service. Just copy-pasting their architecture
 (and choosing hadoop) solutions doesn't fit because we shouldn't **expect**
 such traffic. Our goal is to find balance between scale fees and time to get
 service up and running fast.
 
### Url shortener algorithm choice

### statistics: codebase size & time spent

### possible improvements


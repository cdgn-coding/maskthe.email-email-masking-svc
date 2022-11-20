# Email Masking Service

This service is responsible for masking email addresses. 
It subscribes to a topic of incoming emails and publishes the masked emails to a topic, which is later consumed by the email service.
On the other hand, it manages a database of email addresses and their corresponding masked addresses.

![Architecture](/docs/diagrams/email-masking-svc.png)

## License

    Copyright 2022 Carlos David Gonzalez Nexans
    
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    
      http://www.apache.org/licenses/LICENSE-2.0
    
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
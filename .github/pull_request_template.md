## Overview
<!--  
Provide a high-level description of this change.   
-->

## Functional Demo
<!--  
If applicable, especially if you're adding a new plugin, provide a video or screenshots of how your code works. This will speed up the review process!
-->


## Type of change
<!--  
Check the box below that describes your change best:
--> 

- [ ] Created a new plugin
- [ ] Improved an existing plugin
- [ ] Fixed a bug in an existing plugin
- [ ] Improved contributor utilities or experience

## Related Issue(s)
<!--  
If applicable - add the issue that your PR relates to or closes:
  - use Resolves: #ISSUE_NUMBER to trigger closing of the issue on merge of this PR  
  - use Relates: #ISSUE_NUMBER to indicate relation to an issue, but issue will not close  
-->  

* Resolves: #
* Relates: #

## How To Test
<!--
Provide testing instructions for validating the changes introduced in this PR.
This will serve as a starting point for your reviewers, for functional testing.

If you created a new plugin, you can add a command here which can be used to test authentication.
For example, for the AWS CLI:
  aws s3 ls
-->



## Changelog
<!--  
A one line sentence describing the change that this PR introduces. 
If this has impact over the user experience, your changelog will be included in the release notes of the next stable version of 1Password CLI.

Here are a few guidelines for writing a good changelog:
- Keep your description to a single sentence if you can, and use proper capitalization and punctuation, including a final period.
- Don't use emoji in your description.
- Avoid starting your sentence with "improved" or "fixed". Instead, describe the improvement or say what you fixed.
- Avoid using terminology like "Users are shown" or "You can now" and instead focus on the thing that was changed.

A few examples:

Authenticate the AWS CLI using Touch ID and other unlock options with 1Password Shell Plugins.
The AWS plugin can now be correctly initialized with a default credential, using `op plugin init`.
The AWS plugin now checks for the `AWS_SHARED_CREDENTIALS_FILE` environment variable and attempts to import credentials using the specified file.

For more examples, have a look over 1Password CLI's past release notes: 
https://app-updates.agilebits.com/product_history/CLI2
-->  



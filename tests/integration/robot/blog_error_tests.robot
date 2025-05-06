*** Settings ***
Documentation     Test suite for Blog API error handling
Resource          common.resource
Suite Setup       Setup Test Suite
Suite Teardown    Teardown Test Suite

*** Test Cases ***
Get Non-Existent Blog Post
    ${non_existent_id}=    Set Variable    00000000-0000-0000-0000-000000000000
    Run Keyword And Expect Error    *    Get Blog Post    ${non_existent_id}

Delete Non-Existent Blog Post
    ${non_existent_id}=    Set Variable    00000000-0000-0000-0000-000000000000
    Run Keyword And Expect Error    *    Delete Blog Post    ${non_existent_id}

Update Non-Existent Blog Post
    ${non_existent_id}=    Set Variable    00000000-0000-0000-0000-000000000000
    ${new_title}=    Generate Random String    15
    Run Keyword And Expect Error    *    Update Blog Post    ${non_existent_id}    title=${new_title}

Add Comment To Non-Existent Blog Post
    ${non_existent_id}=    Set Variable    00000000-0000-0000-0000-000000000000
    ${comment_content}=    Generate Random String    30
    ${comment_author}=    Generate Random String    10
    Run Keyword And Expect Error    *    Add Comment To Blog Post    ${non_existent_id}    ${comment_content}    ${comment_author}

Create Blog Post With Empty Title
    ${content}=    Generate Random String    50
    Run Keyword And Expect Error    *    Create Blog Post    ${EMPTY}    ${content}

Create Blog Post With Empty Content
    ${title}=    Generate Random String    15
    Run Keyword And Expect Error    *    Create Blog Post    ${title}    ${EMPTY}

Add Comment With Empty Content
    # Create a blog post first
    ${title}=    Generate Random String    15
    ${content}=    Generate Random String    50
    ${resp}=    Create Blog Post    ${title}    ${content}
    ${blog_id}=    Set Variable    ${resp}[id][value]
    
    # Try to add a comment with empty content
    ${comment_author}=    Generate Random String    10
    Run Keyword And Expect Error    *    Add Comment To Blog Post    ${blog_id}    ${EMPTY}    ${comment_author}
    
    # Clean up
    Delete Blog Post    ${blog_id}

Add Comment With Empty Author
    # Create a blog post first
    ${title}=    Generate Random String    15
    ${content}=    Generate Random String    50
    ${resp}=    Create Blog Post    ${title}    ${content}
    ${blog_id}=    Set Variable    ${resp}[id][value]
    
    # Try to add a comment with empty author
    ${comment_content}=    Generate Random String    30
    Run Keyword And Expect Error    *    Add Comment To Blog Post    ${blog_id}    ${comment_content}    ${EMPTY}
    
    # Clean up
    Delete Blog Post    ${blog_id}
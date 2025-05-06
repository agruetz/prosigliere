*** Settings ***
Documentation     Test suite for Blog API CRUD operations
Resource          common.resource
Suite Setup       Setup Test Suite
Suite Teardown    Teardown Test Suite

*** Test Cases ***
Create Blog Post
    ${title}=    Generate Random String    15
    ${content}=    Generate Random String    50
    ${resp}=    Create Blog Post    ${title}    ${content}
    Should Not Be Empty    ${resp}[id][value]
    Set Suite Variable    ${BLOG_ID}    ${resp}[id][value]

Get Blog Post
    ${resp}=    Get Blog Post    ${BLOG_ID}
    Should Not Be Empty    ${resp}[blog]
    Should Not Be Empty    ${resp}[blog][id][value]
    Should Be Equal    ${resp}[blog][id][value]    ${BLOG_ID}

Update Blog Post
    ${new_title}=    Generate Random String    15
    ${new_content}=    Generate Random String    50
    ${resp}=    Update Blog Post    ${BLOG_ID}    title=${new_title}    content=${new_content}
    
    # Verify the update was successful by getting the post
    ${get_resp}=    Get Blog Post    ${BLOG_ID}
    Should Be Equal    ${get_resp}[blog][title]    ${new_title}
    Should Be Equal    ${get_resp}[blog][content]    ${new_content}

Update Blog Post Title Only
    ${original_resp}=    Get Blog Post    ${BLOG_ID}
    ${original_content}=    Set Variable    ${original_resp}[blog][content]
    
    ${new_title}=    Generate Random String    15
    ${resp}=    Update Blog Post    ${BLOG_ID}    title=${new_title}
    
    # Verify the update was successful by getting the post
    ${get_resp}=    Get Blog Post    ${BLOG_ID}
    Should Be Equal    ${get_resp}[blog][title]    ${new_title}
    Should Be Equal    ${get_resp}[blog][content]    ${original_content}

Update Blog Post Content Only
    ${original_resp}=    Get Blog Post    ${BLOG_ID}
    ${original_title}=    Set Variable    ${original_resp}[blog][title]
    
    ${new_content}=    Generate Random String    50
    ${resp}=    Update Blog Post    ${BLOG_ID}    content=${new_content}
    
    # Verify the update was successful by getting the post
    ${get_resp}=    Get Blog Post    ${BLOG_ID}
    Should Be Equal    ${get_resp}[blog][title]    ${original_title}
    Should Be Equal    ${get_resp}[blog][content]    ${new_content}

Delete Blog Post
    ${resp}=    Delete Blog Post    ${BLOG_ID}
    
    # Verify the post was deleted by trying to get it (should fail)
    Run Keyword And Expect Error    *    Get Blog Post    ${BLOG_ID}
*** Settings ***
Documentation     Test suite for Blog API comment operations
Resource          common.resource
Suite Setup       Setup Test Suite
Suite Teardown    Teardown Test Suite

*** Test Cases ***
Add Comment To Blog Post
    # Create a blog post first
    ${title}=    Generate Random String    15
    ${content}=    Generate Random String    50
    ${resp}=    Create Blog Post    ${title}    ${content}
    Set Test Variable    ${BLOG_ID}    ${resp}[id][value]
    
    # Add a comment to the blog post
    ${comment_content}=    Generate Random String    30
    ${comment_author}=    Generate Random String    10
    ${comment_resp}=    Add Comment To Blog Post    ${BLOG_ID}    ${comment_content}    ${comment_author}
    
    # Verify the comment was added by getting the blog post
    ${get_resp}=    Get Blog Post    ${BLOG_ID}
    Should Not Be Empty    ${get_resp}[blog][comments]
    Length Should Be    ${get_resp}[blog][comments]    1
    Should Be Equal    ${get_resp}[blog][comments][0][content]    ${comment_content}
    Should Be Equal    ${get_resp}[blog][comments][0][author]    ${comment_author}

Add Multiple Comments To Blog Post
    # Add more comments to the same blog post
    ${comment_content1}=    Generate Random String    30
    ${comment_author1}=    Generate Random String    10
    Add Comment To Blog Post    ${BLOG_ID}    ${comment_content1}    ${comment_author1}
    
    ${comment_content2}=    Generate Random String    30
    ${comment_author2}=    Generate Random String    10
    Add Comment To Blog Post    ${BLOG_ID}    ${comment_content2}    ${comment_author2}
    
    # Verify all comments are present
    ${get_resp}=    Get Blog Post    ${BLOG_ID}
    Should Not Be Empty    ${get_resp}[blog][comments]
    Length Should Be    ${get_resp}[blog][comments]    3
    
    # Verify comment count in list view
    ${list_resp}=    List Blog Posts
    ${found_blog}=    Set Variable    ${NONE}
    FOR    ${blog}    IN    @{list_resp}[blogs]
        Run Keyword If    '${blog}[id][value]' == '${BLOG_ID}'    Set Test Variable    ${found_blog}    ${blog}
    END
    
    Should Not Be Equal    ${found_blog}    ${NONE}
    Should Be Equal As Integers    ${found_blog}[commentCount]    3

Cleanup Test Blog Post
    # Clean up the blog post we created
    Delete Blog Post    ${BLOG_ID}
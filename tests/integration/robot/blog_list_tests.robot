*** Settings ***
Documentation     Test suite for Blog API listing operations
Resource          common.resource
Suite Setup       Setup Test Suite
Suite Teardown    Teardown Test Suite

*** Test Cases ***
List Blog Posts
    # Create a few blog posts first
    ${title1}=    Generate Random String    15
    ${content1}=    Generate Random String    50
    ${resp1}=    Create Blog Post    ${title1}    ${content1}
    Set Test Variable    ${BLOG_ID_1}    ${resp1}[id][value]
    
    ${title2}=    Generate Random String    15
    ${content2}=    Generate Random String    50
    ${resp2}=    Create Blog Post    ${title2}    ${content2}
    Set Test Variable    ${BLOG_ID_2}    ${resp2}[id][value]
    
    ${title3}=    Generate Random String    15
    ${content3}=    Generate Random String    50
    ${resp3}=    Create Blog Post    ${title3}    ${content3}
    Set Test Variable    ${BLOG_ID_3}    ${resp3}[id][value]
    
    # List all blog posts
    ${list_resp}=    List Blog Posts
    Should Not Be Empty    ${list_resp}[blogs]
    
    # Verify our created posts are in the list
    ${found_blog1}=    Set Variable    ${FALSE}
    ${found_blog2}=    Set Variable    ${FALSE}
    ${found_blog3}=    Set Variable    ${FALSE}
    
    FOR    ${blog}    IN    @{list_resp}[blogs]
        Run Keyword If    '${blog}[id][value]' == '${BLOG_ID_1}'    Set Test Variable    ${found_blog1}    ${TRUE}
        Run Keyword If    '${blog}[id][value]' == '${BLOG_ID_2}'    Set Test Variable    ${found_blog2}    ${TRUE}
        Run Keyword If    '${blog}[id][value]' == '${BLOG_ID_3}'    Set Test Variable    ${found_blog3}    ${TRUE}
    END
    
    Should Be True    ${found_blog1}
    Should Be True    ${found_blog2}
    Should Be True    ${found_blog3}

List Blog Posts With Pagination
    # Create more blog posts to ensure we have enough for pagination
    FOR    ${i}    IN RANGE    5
        ${title}=    Generate Random String    15
        ${content}=    Generate Random String    50
        Create Blog Post    ${title}    ${content}
    END
    
    # Get first page with 2 items
    ${page1_resp}=    List Blog Posts    page_size=2
    Should Not Be Empty    ${page1_resp}[blogs]
    Length Should Be    ${page1_resp}[blogs]    2
    Should Not Be Empty    ${page1_resp}[nextPageToken]
    
    # Get second page
    ${page2_resp}=    List Blog Posts    page_size=2    page_token=${page1_resp}[nextPageToken]
    Should Not Be Empty    ${page2_resp}[blogs]
    Length Should Be    ${page2_resp}[blogs]    2
    
    # Verify different posts on different pages
    ${page1_ids}=    Create List
    FOR    ${blog}    IN    @{page1_resp}[blogs]
        Append To List    ${page1_ids}    ${blog}[id][value]
    END
    
    FOR    ${blog}    IN    @{page2_resp}[blogs]
        Should Not Contain    ${page1_ids}    ${blog}[id][value]
    END

Cleanup Test Blog Posts
    # Clean up the blog posts we created
    Run Keyword If Defined    ${BLOG_ID_1}    Delete Blog Post    ${BLOG_ID_1}
    Run Keyword If Defined    ${BLOG_ID_2}    Delete Blog Post    ${BLOG_ID_2}
    Run Keyword If Defined    ${BLOG_ID_3}    Delete Blog Post    ${BLOG_ID_3}
*** Settings ***
Documentation     Common resources for Blog API testing
Library           RequestsLibrary
Library           Collections
Library           String
Library           OperatingSystem
Library           DateTime

*** Variables ***
${BASE_URL}       http://localhost:8080
${API_PATH}       /v1/posts
${CONTENT_TYPE}   application/json

*** Keywords ***
Setup Test Suite
    Create Session    blog_api    ${BASE_URL}    verify=True

Teardown Test Suite
    Delete All Sessions

Generate Random String
    [Arguments]    ${length}=10
    ${random_string}=    String.Generate Random String    ${length}    [LETTERS]
    [Return]    ${random_string}

Create Blog Post
    [Arguments]    ${title}    ${content}
    ${body}=    Create Dictionary    title=${title}    content=${content}
    ${resp}=    POST On Session    blog_api    ${API_PATH}    json=${body}    expected_status=200
    [Return]    ${resp.json()}

Get Blog Post
    [Arguments]    ${post_id}
    ${resp}=    GET On Session    blog_api    ${API_PATH}/${post_id}    expected_status=200
    [Return]    ${resp.json()}

Delete Blog Post
    [Arguments]    ${post_id}
    ${resp}=    DELETE On Session    blog_api    ${API_PATH}/${post_id}    expected_status=200
    [Return]    ${resp}

Update Blog Post
    [Arguments]    ${post_id}    ${title}=${EMPTY}    ${content}=${EMPTY}
    ${body}=    Create Dictionary
    Run Keyword If    '${title}' != '${EMPTY}'    Set To Dictionary    ${body}    title=${title}
    Run Keyword If    '${content}' != '${EMPTY}'    Set To Dictionary    ${body}    content=${content}
    ${resp}=    PATCH On Session    blog_api    ${API_PATH}/${post_id}    json=${body}    expected_status=200
    [Return]    ${resp}

Add Comment To Blog Post
    [Arguments]    ${post_id}    ${content}    ${author}
    ${body}=    Create Dictionary    content=${content}    author=${author}
    ${resp}=    POST On Session    blog_api    ${API_PATH}/${post_id}/comment    json=${body}    expected_status=200
    [Return]    ${resp}

List Blog Posts
    [Arguments]    ${page_size}=${EMPTY}    ${page_token}=${EMPTY}
    ${params}=    Create Dictionary
    Run Keyword If    '${page_size}' != '${EMPTY}'    Set To Dictionary    ${params}    pageSize=${page_size}
    Run Keyword If    '${page_token}' != '${EMPTY}'    Set To Dictionary    ${params}    pageToken=${page_token}
    ${resp}=    GET On Session    blog_api    ${API_PATH}    params=${params}    expected_status=200
    [Return]    ${resp.json()}
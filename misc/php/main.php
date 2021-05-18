<?php

function PPP2021_send_post_request($data)
{
    $args      = array(
        'body' => json_encode($data),
        'timeout' => '30',
        'blocking' => true,
        'headers' => array(
            'Content-Type: application/json'
        )
    );
    $response  = wp_remote_post(<backend-url>, $args);
    $http_code = wp_remote_retrieve_response_code($response);
    if ($http_code !== 200) {
        error_log('error: network error posting to verifiable datastructure');
        error_log(json_encode($response));
    }
}

function PPP2021_commit_post_transition($new_status, $old_status, $post_id)
{

    if (is_null($post_id)) {
        error_log('error: $post_id is null, unable to commit post transition to verifiable datastructure');
        return;
    }

    $post = get_post($post_id); // example below
    if (is_null($post)) {
        error_log('error: $post is null, unable to commit post transition to verifiable datastructure');
        return;
    }

    if ($post->post_author) {
        // include display_name in the verifiable datastructure
        $post->post_author_display_name = get_the_author_meta('display_name', $post->post_author);
    }

    $json_data = array(
        'ID' => $post->ID,
        'Data' => PPP2021_get_post_hash($post_id)
    );
    PPP2021_send_post_request($json_data);
}

add_action('transition_post_status', 'PPP2021_commit_post_transition', 10, 3);

function PPP2021_get_post_hash($post_id)
{
	$post = get_post($post_id);
    if ($post->post_author) {
        // include display_name in the verifiable datastructure
        $post->post_author_display_name = get_the_author_meta('display_name', $post->post_author);
    }
	return hash('sha256', json_encode($post));
}


add_action('transition_post_status', 'PPP2021_commit_post_transition', 10, 3);

/*
Example value of $post:
{
"ID": 5,
"post_author": "1",
"post_author_display_name": "newsadmin",
"post_date": "2021-01-08 02:22:08",
"post_date_gmt": "2021-01-08 02:22:08",
"post_content": "<!-- wp:paragraph -->\n<p>Every month, the archival institutions of this nation unleash tiny particles of the past in a frenzy of online revelry.</p>\n<!-- /wp:paragraph -->\n\n<!-- wp:paragraph -->\n<p>Is there room in your mind for unuseful details? Or, would you make some, despite long odds the material could one day prove practicable?</p>\n<!-- /wp:paragraph -->",
"post_title": "The Record Keepers’ Rave",
"post_excerpt": "",
"post_status": "publish",
"comment_status": "open",
"ping_status": "open",
"post_password": "",
"post_name": "the-record-keepers-rave",
"to_ping": "",
"pinged": "",
"post_modified": "2021-01-08 02:48:54",
"post_modified_gmt": "2021-01-08 02:48:54",
"post_content_filtered": "",
"post_parent": 0,
"guid": "http://newsprovenance.kinsta.cloud/?p=5",
"menu_order": 0,
"post_type": "post",
"post_mime_type": "",
"comment_count": "0",
"filter": "raw"
}
*/
?>

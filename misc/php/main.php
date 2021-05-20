<?php
/**
 * @package Hello_Dolly
 * @version 1.7.2
 */
/*
Plugin Name: Hello Dolly
Plugin URI: http://wordpress.org/plugins/hello-dolly/
Description: This is not just a plugin, it symbolizes the hope and enthusiasm of an entire generation summed up in two words sung most famously by Louis Armstrong: Hello, Dolly. When activated you will randomly see a lyric from <cite>Hello, Dolly</cite> in the upper right of your admin screen on every page.
Author: Matt Mullenweg
Version: 1.7.2
Author URI: http://ma.tt/
*/

function PPP2021_send_post_request($url, $data)
{
    $args = array(
        'body' => json_encode($data) ,
        'timeout' => '30',
        'blocking' => true,
        'headers' => array(
            'Content-Type: application/json'
        )
    );
    $response = wp_remote_post($url, $args);
    $http_code = wp_remote_retrieve_response_code($response);
    if ($http_code !== 200)
    {
        error_log('error: network error posting to verifiable datastructure');
        error_log(json_encode($response));
    }
    return wp_remote_retrieve_body($response);
}

function PPP2021_commit_post_transition($new_status, $old_status, $post_id)
{

    // error check, can ignore
    if (is_null($post_id))
    {
        error_log('error: $post_id is null, unable to commit post transition to verifiable datastructure');
        return;
    }

    // get some of the basic information about the post, from the ID
    $post = get_post($post_id); // example below
    if (is_null($post))
    {
        error_log('error: $post is null, unable to commit post transition to verifiable datastructure');
        return;
    }

    // get some additional information about the post, extra can be explored at end
    if ($post->post_author)
    {
        // include display_name in the verifiable datastructure
        $post->post_author_display_name = get_the_author_meta('display_name', $post->post_author);
    }

    // this is what will get HTTP posted to log/ map server
    // it corresponds to WordpressPost type
    // https://github.com/z-tech/blue/blob/main/src/types/wordpressPost.go#L5
    $data = array(
        'ID' => $post->ID,
        'Data' => PPP2021_get_post_hash($post_id)
    );
    PPP2021_send_post_request('http://ec2-54-210-116-133.compute-1.amazonaws.com:8080/v1/commitWordpressPost', $data);
}

add_action('transition_post_status', 'PPP2021_commit_post_transition', 10, 3);

function PPP2021_get_post_hash($post_id)
{
    // same as lines 42-53
    $post = get_post($post_id);
    if ($post->post_author)
    {
        // include display_name in the verifiable datastructure
        $post->post_author_display_name = get_the_author_meta('display_name', $post->post_author);
    }
    // runs that object through sha256 cryptographic hash
    // THIS IS WHERE YOU WANNA MAKE YOUR CHANGES
    // $stringified_schema_org_format = '<xml><ID>' . $post->ID . '</ID></xml>';
    // return hash('sha256', $stringified_schema_org_format);
    return base64_encode(hash('sha256', json_encode($post), true)); // true means binary instead of lowercase hexits
}

function PPP2021_get_post_log_hash($post_id)
{
    // same as lines 42-53
    $post = get_post($post_id);
    if ($post->post_author)
    {
        // include display_name in the verifiable datastructure
        $post->post_author_display_name = get_the_author_meta('display_name', $post->post_author);
    }
	$RFC6962LeafHashPrefix = chr(0); // https://www.php.net/manual/en/function.chr.php
	$log_leaf_string = $RFC6962LeafHashPrefix . $post->ID . ',' . PPP2021_get_post_hash($post_id);
    return base64_encode(hash('sha256', $log_leaf_string, true)); // true means binary instead of lowercase hexits
}

function PPP2021_get_proofs($post_id)
{
    $hash = PPP2021_get_post_hash($post_id);
    $data = array(
        'ID' => $post_id,
        'Data' => $hash
    );
    $result = PPP2021_send_post_request('server-url', $data);
    return json_decode($result);
}

function PPP2021_get_tree_roots()
{
	$result = wp_remote_get('server-url');
	$result = wp_remote_retrieve_body($result);
    return json_decode($result);
}

function PPP2021_get_pretty_printed_proofs($post_id)
{
    return json_encode(PPP2021_get_proofs($post_id) , JSON_PRETTY_PRINT);
}

function PPP2021_get_pretty_printed_tree_roots()
{
    return json_encode(PPP2021_get_tree_roots() , JSON_PRETTY_PRINT);
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

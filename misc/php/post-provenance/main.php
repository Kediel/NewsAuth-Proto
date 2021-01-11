<?php
/*
Contributors: alz236
Description: A plugin to commit post content to a verifiable datastructure
License: GPLv2 or later
Plugin Name: Post Provenance Project
Plugin URI: https://github.com/z-tech/blue/tree/main/misc/php
Requires at least: 4.6
Requires PHP: 7.3
Stable tag: 0.0.1
Tested up to: 5.6
Version: 0.0.1
*/

defined('ABSPATH') or die('No peeking at this page!');

// Register the menu
add_action("admin_menu", "PPP2021_plugin_menu_func");
function PPP2021_plugin_menu_func()
{
    add_submenu_page("options-general.php", // Which menu parent
        "Post Provenance", // Page title
        "Post Provenance", // Menu title
        "manage_options", // Minimum capability (manage_options is an easy way to target Admins)
        "postprovenance", // Menu slug
        "PPP2021_plugin_options" // Callback that prints the markup
        );
}

// Print the markup for the page
function PPP2021_plugin_options()
{
    if (!current_user_can("manage_options")) {
        wp_die(__("You do not have sufficient permissions to access this page."));
    }

    if (isset($_GET['status']) && $_GET['status'] == 'success') {
?>
     <div id="message" class="updated notice is-dismissible">
          <p><?php
        _e("Settings updated!", "provenance-api");
?></p>
          <button type="button" class="notice-dismiss">
              <span class="screen-reader-text"><?php
        _e("Dismiss this notice.", "provenance-api");
?></span>
          </button>
      </div>
    <?php
    }

?>
 <form method="post" action="<?php
    echo admin_url('admin-post.php');
?>">

  <input type="hidden" name="action" value="provenance_credentials_submit" />

  <h3>Provenance Credentials</h3>

  <p>
    <label><?php
    _e("Provenance API url:", "provenance-api");
?></label>
    <input class="" type="text" name="provenance_api_url" value="<?php
    echo get_option('provenance_api_url');
?>" />
  </p>
  <p>
    <label><?php
    _e("Provenance API key:", "provenance-api");
?></label>
    <input class="" type="password" name="provenance_api_key" value="<?php
    echo get_option('provenance_api_key');
?>" />
  </p>

  <input class="button button-primary" type="submit" value="<?php
    _e("Save", "provenance-api");
?>" />

  </form>

<?php

}

add_action('admin_post_provenance_credentials_submit', 'PPP2021_credentials_handle_save');

function PPP2021_credentials_handle_save()
{

    // Get the options that were sent
    $provenance_url  = (!empty($_POST["provenance_api_url"])) ? $_POST["provenance_api_url"] : NULL;
    $provenance_key = (!empty($_POST["provenance_api_key"])) ? $_POST["provenance_api_key"] : NULL;

    // Validation would go here

    // Update the values
    update_option("PPP2021_provenance_api_url", $provenance_url, TRUE);
    update_option("PPP2021_provenance_api_key", $provenance_key, TRUE);

    // Redirect back to settings page
    // The ?page=github corresponds to the "slug"
    // set in the fourth parameter of add_submenu_page() above.
    $redirect_url = get_bloginfo("url") . "/wp-admin/options-general.php?page=postprovenance&status=success";
    header("Location: " . $redirect_url);
    exit;
}

function PPP2021_send_post_request($url, $api_token, $data)
{
    $args      = array(
        'body' => $data,
        'timeout' => '30',
        'blocking' => true,
        'headers' => array(
            'Content-Type: application/json',
            'Authorization: Bearer ' . $api_token
        )
    );
    $response  = wp_remote_post($url . '/v1/proveWordpressPost', $args);
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

	$api_adress = get_option('PPP2021_provenance_api_url');
	$api_key = get_option('PPP2021_provenance_api_key');
    if (is_null($api_adress) || is_null($api_key)) {
        return;
    }

    $json_data = array(
        'ID' => $post->ID,
        'Data' => json_encode($post)
    );
    PPP2021_send_post_request($api_adress, $api_key, $json_data);
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

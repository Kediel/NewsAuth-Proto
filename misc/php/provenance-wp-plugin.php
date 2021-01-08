<?php
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

function send_post_request($url, $token, $data)
{
  $ch = curl_init($url);
  $jsonDataEncoded = json_encode($data);
  curl_setopt($ch, CURLOPT_POST, 1);
  curl_setopt($ch, CURLOPT_POSTFIELDS, $jsonDataEncoded);
  curl_setopt($ch, CURLOPT_HTTPHEADER, array(
    'Content-Type: application/json',
    'Authorization: Bearer ' . $token
  ));
  $result = curl_exec($ch);
  if ($result === false)
  {
    error_log('error: network error posting to verifiable datastructure');
    error_log(json_encode($result));
  }
}

function commit_post_transition($new_status, $old_status, $post_id)
{
  $VD_API_ADDRESS = 'FILL_THIS_VALUE';
  $VD_API_KEY = 'FILL_THIS_VALUE';

  if (is_null($post_id))
  {
    error_log('error: $post_id is null, unable to commit post transition to verifiable datastructure');
    return;
  }

  $post = get_post($post_id); // example above
  if (is_null($post))
  {
    error_log('error: $post is null, unable to commit post transition to verifiable datastructure');
    return;
  }

  if ($post->post_author)
  {
    // include display_name in the verifiable datastructure
    $post->post_author_display_name = get_the_author_meta('display_name', $post->post_author);
  }

  if (base64_encode($VD_API_ADDRESS) === 'RklMTF9USElTX1ZBTFVF' || base64_encode($VD_API_KEY) === 'RklMTF9USElTX1ZBTFVF')
  {
    error_log('error: both $VD_API_ADDRESS and $VD_API_KEY must be set');
    return;
  }

  $json_data = array(
    'PostID' => $post->ID,
    'Post' => $post
  );
  send_post_request($VD_API_ADDRESS, $VD_API_KEY, $json_data);
}

add_action('transition_post_status', 'commit_post_transition', 10, 3);
?>

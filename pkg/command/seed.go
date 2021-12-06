package command

import (
	"context"

	"github.com/neoxelox/odin/internal"
	"github.com/neoxelox/odin/internal/class"
	"github.com/neoxelox/odin/internal/core"
	"github.com/neoxelox/odin/internal/database"
)

const communities = `
INSERT INTO "community" ("id", "address", "name", "categories", "pinned_ids")
VALUES ('9bsv0s5a5rsg02purd40', 'Rambla les ferreries nº44', 'Rambla les ferreries nº44', '{Suministros, Desagües, Cerrajería, Ascensor, Estructural, Zonas Comunes, Otros}', '{}')
;
`

const users = `
INSERT INTO "user" ("id", "phone", "name", "email", "picture", "birthday", "language", "is_banned")
VALUES ('9bsv0s2l1bfg034l8so0', '+34722560561', 'Susana García', '', 'https://images.unsplash.com/photo-1619895862022-09114b41f16f?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Mnx8cHJvZmlsZSUyMHBpY3R1cmV8ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8sog', '+34637246690', 'Alicia Navarro', '', 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8M3x8cHJvZmlsZSUyMHBpY3R1cmV8ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8sp0', '+34730891342', 'Mauricio Fernandez', '', 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Nnx8cHJvZmlsZSUyMHBpY3R1cmV8ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8spg', '+34715606309', 'Lourdes Gonzalez', '', 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8NXx8cHJvZmlsZSUyMHBpY3R1cmV8ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8sq0', '+34603118899', 'Juana Sanchez', '', 'https://images.unsplash.com/photo-1508214751196-bcfd4ca60f91?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8OXx8cHJvZmlsZSUyMHBpY3R1cmV8ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8sqg', '+34739958395', 'Daniel Garcia', '', 'https://images.unsplash.com/photo-1544723795-3fb6469f5b39?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8OHx8cHJvZmlsZSUyMHBpY3R1cmV8ZW58MHx8MHx8&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8sr0', '+34760038121', 'Norma Benitez', '', 'https://images.unsplash.com/photo-1628890923662-2cb23c2e0cfe?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MTZ8fHByb2ZpbGUlMjBwaWN0dXJlfGVufDB8fDB8fA%3D%3D&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8srg', '+34628263861', 'Alicia Torres', '', 'https://images.unsplash.com/photo-1521252659862-eec69941b071?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MTh8fHByb2ZpbGUlMjBwaWN0dXJlfGVufDB8fDB8fA%3D%3D&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8ss0', '+34791927787', 'Asunción Flores', '', 'https://media.istockphoto.com/photos/one-beautiful-woman-looking-at-the-camera-in-profile-picture-id1303539316?b=1&k=20&m=1303539316&s=170667a&w=0&h=ePGGvjsOR__-R2KSvZ67xXl2x-CkVzKg8q_WtvqLww0=', current_date, '', false),
       ('9bsv0s2l1bfg034l8ssg', '+34635418380', 'Sergio Ruiz', '', 'https://media.istockphoto.com/photos/smiling-man-with-hat-and-sunglasses-picture-id953079238?b=1&k=20&m=953079238&s=170667a&w=0&h=nV7ep-uPYLJhRtYlHB3SCEJ1Sksd-P5o1-sZDslffJI=', current_date, '', false),
       ('9bsv0s2l1bfg034l8st0', '+34642172982', 'María Romero', '', 'https://media.istockphoto.com/photos/pleasant-young-indian-woman-freelancer-consult-client-via-video-call-picture-id1300972573?b=1&k=20&m=1300972573&s=170667a&w=0&h=xuAsEkMkoBbc5Nh-nButyq3DU297V_tnak-60VarrR0=', current_date, '', false),
       ('9bsv0s2l1bfg034l8stg', '+34790564031', 'Víctor Benitez', '', 'https://media.istockphoto.com/photos/millennial-male-team-leader-organize-virtual-workshop-with-employees-picture-id1300972574?b=1&k=20&m=1300972574&s=170667a&w=0&h=2nBGC7tr0kWIU8zRQ3dMg-C5JLo9H2sNUuDjQ5mlYfo=', current_date, '', false),
       ('9bsv0s2l1bfg034l8su0', '+34784206955', 'María Ramirez', '', 'https://media.istockphoto.com/photos/smile-girl-at-beach-picture-id477151294?b=1&k=20&m=477151294&s=170667a&w=0&h=T3n6eYUeg26yjssHbVKhZTbB8reMcBNAWut_ut5u8yY=', current_date, '', false),
       ('9bsv0s2l1bfg034l8sug', '+34740052904', 'Patricia Santos', '', 'https://images.unsplash.com/photo-1625897428517-7e2062829a26?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MjB8fHByb2ZpbGUlMjBwaWN0dXJlfGVufDB8fDB8fA%3D%3D&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8sv0', '+34717563445', 'Luis Alvarez', '', 'https://images.unsplash.com/photo-1547425260-76bcadfb4f2c?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8NHx8cGVyc29ufGVufDB8fDB8fA%3D%3D&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8svg', '+34735690477', 'Juan Garcia', '', 'https://images.unsplash.com/photo-1599566150163-29194dcaad36?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MTh8fHBlcnNvbnxlbnwwfHwwfHw%3D&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8t00', '+34639894510', 'Elena Ramirez', '', 'https://images.unsplash.com/photo-1580489944761-15a19d654956?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Mjh8fHBlcnNvbnxlbnwwfHwwfHw%3D&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8t0g', '+34631226695', 'Cesar Sosa', '', 'https://images.unsplash.com/photo-1540569014015-19a7be504e3a?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8Mzh8fHBlcnNvbnxlbnwwfHwwfHw%3D&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8t10', '+34618407958', 'Ricardo Romero', '', 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8NDF8fHBlcnNvbnxlbnwwfHwwfHw%3D&auto=format&fit=crop&w=500&q=60', current_date, '', false),
       ('9bsv0s2l1bfg034l8t1g', '+34737106464', 'Álex Aguirre', '', 'https://images.unsplash.com/photo-1598411072028-c4642d98352c?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxzZWFyY2h8NDR8fHBlcnNvbnxlbnwwfHwwfHw%3D&auto=format&fit=crop&w=500&q=60', current_date, '', false)
;
`

const memberships = `
INSERT INTO "membership" ("id", "user_id", "community_id", "door", "role")
VALUES ('9bsv0s65sclg02m57a20', '9bsv0s2l1bfg034l8so0', '9bsv0s5a5rsg02purd40', '3º 4ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a2g', '9bsv0s2l1bfg034l8sog', '9bsv0s5a5rsg02purd40', '3º 3ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a30', '9bsv0s2l1bfg034l8sp0', '9bsv0s5a5rsg02purd40', '3º 3ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a3g', '9bsv0s2l1bfg034l8spg', '9bsv0s5a5rsg02purd40', '1º 1ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a40', '9bsv0s2l1bfg034l8sq0', '9bsv0s5a5rsg02purd40', '3º 2ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a4g', '9bsv0s2l1bfg034l8sqg', '9bsv0s5a5rsg02purd40', '2º 4ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a50', '9bsv0s2l1bfg034l8sr0', '9bsv0s5a5rsg02purd40', '1º 2ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a5g', '9bsv0s2l1bfg034l8srg', '9bsv0s5a5rsg02purd40', 'Bajo 2ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a60', '9bsv0s2l1bfg034l8ss0', '9bsv0s5a5rsg02purd40', '2º 1ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a6g', '9bsv0s2l1bfg034l8ssg', '9bsv0s5a5rsg02purd40', 'Bajo 2ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a70', '9bsv0s2l1bfg034l8st0', '9bsv0s5a5rsg02purd40', 'Bajo 2ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a7g', '9bsv0s2l1bfg034l8stg', '9bsv0s5a5rsg02purd40', '1º 2ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a80', '9bsv0s2l1bfg034l8su0', '9bsv0s5a5rsg02purd40', '2º 1ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a8g', '9bsv0s2l1bfg034l8sug', '9bsv0s5a5rsg02purd40', '3º 3ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a90', '9bsv0s2l1bfg034l8sv0', '9bsv0s5a5rsg02purd40', '1º 2ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57a9g', '9bsv0s2l1bfg034l8svg', '9bsv0s5a5rsg02purd40', '1º 4ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57aa0', '9bsv0s2l1bfg034l8t00', '9bsv0s5a5rsg02purd40', '3º 1ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57aag', '9bsv0s2l1bfg034l8t0g', '9bsv0s5a5rsg02purd40', 'Bajo 4ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57ab0', '9bsv0s2l1bfg034l8t10', '9bsv0s5a5rsg02purd40', '1º 3ª', 'RESIDENT'),
       ('9bsv0s65sclg02m57abg', '9bsv0s2l1bfg034l8t1g', '9bsv0s5a5rsg02purd40', '3º 1ª', 'RESIDENT')
;
`

const posts = `
INSERT INTO "post" ("id", "thread_id", "creator_id", "last_history_id", "type", "priority", "recipient_ids", "voter_ids", "created_at")
VALUES ('9bsv0s2u84o002rhqpa0', null, '9bsv0s65sclg02m57a2g', 'TODO', 'PUBLICATION', null, null, '{9bsv0s65sclg02m57a30, 9bsv0s65sclg02m57a3g}', current_timestamp - interval '13 days' - interval '2 hours'),
       ('9bsv0s6k1b9g02sp84c0', '9bsv0s2u84o002rhqpa0', '9bsv0s65sclg02m57a3g', 'TODO', 'PUBLICATION', null, null, '{9bsv0s65sclg02m57a2g}', current_timestamp - interval '13 days' - interval '1 hour'),
       ('9bsv0s2r0u8g02qj5ai0', '9bsv0s2u84o002rhqpa0', '9bsv0s65sclg02m57a40', 'TODO', 'PUBLICATION', null, null, '{9bsv0s65sclg02m57a2g}', current_timestamp - interval '13 days'),
       ('9bsv0s2tchag02t6gg10', null, '9bsv0s65sclg02m57a8g', 'TODO', 'PUBLICATION', null, null, '{}', current_timestamp - interval '11 days' - interval '5 hour'),
       ('9bsv0s09onsg02s5c8ag', '9bsv0s2tchag02t6gg10', '9bsv0s65sclg02m57a7g', 'TODO', 'PUBLICATION', null, null, '{}', current_timestamp - interval '11 days' - interval '3 hour'),
       ('9bsv0s2ccr9g02gl6r50', '9bsv0s2tchag02t6gg10', '9bsv0s65sclg02m57a9g', 'TODO', 'PUBLICATION', null, null, '{}', current_timestamp - interval '11 days' - interval '4 hour' - interval '32 minutes'),
       ('9bsv0s4ilem002n1msv0', null, '9bsv0s65sclg02m57ab0', 'TODO', 'ISSUE', '9', null, '{9bsv0s65sclg02m57a9g, 9bsv0s65sclg02m57a80, 9bsv0s65sclg02m57a60, 9bsv0s65sclg02m57a50}', current_timestamp - interval '9 day' + interval '3 hours'),
       ('9bsv0s4ilem002n1msvg', '9bsv0s4ilem002n1msv0', '9bsv0s65sclg02m57abg', 'TODO', 'PUBLICATION', null, null, '{}', current_timestamp - interval '9 days' + interval '5 hour' - interval '47 minutes'),
       ('9bsv0s4ilem002n1mt00', null, '9bsv0s65sclg02m57a40', 'TODO', 'PUBLICATION', null, null, '{}', current_timestamp - interval '8 days' - interval '7 hour' - interval '02 minutes'),
       ('9bsv0s4ilem002n1mt0g', null, '9bsv0s65sclg02m57ab0', 'TODO', 'PUBLICATION', null, null, '{}', current_timestamp - interval '8 days' - interval '1 hour' - interval '05 minutes'),
       ('9bsv0s4ilem002n1mt10', null, '9bsv0s65sclg02m57a2g', 'TODO', 'ISSUE', '7', null, '{9bsv0s65sclg02m57a80, 9bsv0s65sclg02m57a50}', current_timestamp - interval '6 day' - interval '2 hours' - interval '56 minutes'),
       ('9bsv0s4ilem002n1mt1g', '9bsv0s4ilem002n1mt10', '9bsv0s65sclg02m57a50', 'TODO', 'PUBLICATION', null, null, '{}', current_timestamp - interval '6 days' - interval '1 hour' - interval '13 minutes'),
       ('9bsv0s4ilem002n1mt20', null, '9bsv0s65sclg02m57a4g', 'TODO', 'PUBLICATION', null, null, '{9bsv0s65sclg02m57a60, 9bsv0s65sclg02m57a6g, 9bsv0s65sclg02m57a70, 9bsv0s65sclg02m57a7g, 9bsv0s65sclg02m57a80, 9bsv0s65sclg02m57a8g, 9bsv0s65sclg02m57a90, 9bsv0s65sclg02m57a9g, 9bsv0s65sclg02m57aa0, 9bsv0s65sclg02m57aag}', current_timestamp - interval '2 days' - interval '7 hour' - interval '41 minutes'),
       ('9bsv0s4ilem002n1mt2g', '9bsv0s4ilem002n1mt20', '9bsv0s65sclg02m57a6g', 'TODO', 'PUBLICATION', null, null, '{9bsv0s65sclg02m57a4g}', current_timestamp - interval '6 days' + interval '1 hour' - interval '13 minutes'),
       ('9bsv0s4ilem002n1mt30', '9bsv0s4ilem002n1mt2g', '9bsv0s65sclg02m57a4g', 'TODO', 'PUBLICATION', null, null, '{}', current_timestamp - interval '6 days' + interval '1 hour' - interval '26 minutes'),
       ('9bsv0s4ilem002n1mt3g', '9bsv0s4ilem002n1mt20', '9bsv0s65sclg02m57a9g', 'TODO', 'PUBLICATION', null, null, '{}', current_timestamp - interval '6 days' + interval '5 hour' - interval '49 minutes')
;
`

const post_histories = `
INSERT INTO "post_history" ("id", "post_id", "updator_id", "message", "categories", "state", "media", "widgets", "created_at")
VALUES ('9bsv0s4ilem002n1mt40', '9bsv0s2u84o002rhqpa0', '9bsv0s65sclg02m57a2g', '¿A alguien le sobran huevos? No me acordaba de que hoy era festivo.', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '13 days' - interval '2 hours'),
       ('9bsv0s4ilem002n1mt4g', '9bsv0s6k1b9g02sp84c0', '9bsv0s65sclg02m57a3g', 'Puedes pasarte por mi piso si quieres y te doy un par.', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '13 days' - interval '1 hour'),
       ('9bsv0s4ilem002n1mt50', '9bsv0s2r0u8g02qj5ai0', '9bsv0s65sclg02m57a40', 'A mí me ha pasado lo mismo 😂🤣', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '13 days'),
       ('9bsv0s4ilem002n1mt5g', '9bsv0s2tchag02t6gg10', '9bsv0s65sclg02m57a8g', 'Me acabo de quedar sin coche. Alguien sabe qué bus tengo que coger para ir hasta la estación de trenes?', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '11 days' - interval '5 hour'),
       ('9bsv0s4ilem002n1mt60', '9bsv0s09onsg02s5c8ag', '9bsv0s65sclg02m57a7g', 'Sí, puedes coger la línea 7 que está cruzando la calle, frente al supermercado. El bus pasa cada 5 minutos más o menos así que no deberías tener ningún problema.', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '11 days' - interval '3 hour'),
       ('9bsv0s4ilem002n1mt6g', '9bsv0s2ccr9g02gl6r50', '9bsv0s65sclg02m57a9g', 'Yo siempre cojo la línea 7.', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '11 days' - interval '4 hour' - interval '32 minutes'),
       ('9bsv0s4ilem002n1mt70', '9bsv0s4ilem002n1msv0', '9bsv0s65sclg02m57ab0', 'Buenos días. ¿A alguien más no le funciona la tele? No sintoniza ningún canal. Creo que la antena de la terraza se rompió con la tormenta ayer...', '{}', 'PENDING', '{}', '{"Poll": null}', current_timestamp - interval '9 day' + interval '3 hours'),
       ('9bsv0s4ilem002n1mt7g', '9bsv0s4ilem002n1msvg', '9bsv0s65sclg02m57abg', 'Sí, no funciona', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '9 days' + interval '5 hour' - interval '47 minutes'),
       ('9bsv0s4ilem002n1mt80', '9bsv0s4ilem002n1mt00', '9bsv0s65sclg02m57a40', 'Mañana de 4 a 7 de la tarde celebraremos el cumpleaños de mi hijo en mi piso. Espero que no os importe que pongamos música durante la pequeña fiesta.', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '8 days' - interval '7 hour' - interval '02 minutes'),
       ('9bsv0s4ilem002n1mt8g', '9bsv0s4ilem002n1mt0g', '9bsv0s65sclg02m57ab0', 'He hecho asado argentino y me ha sobrado muchísimo, si alguien quiere que me lo diga!', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '8 days' - interval '1 hour' - interval '05 minutes'),
       ('9bsv0s4ilem002n1mt90', '9bsv0s4ilem002n1mt10', '9bsv0s65sclg02m57a2g', 'El ascensor parece que no funciona. Cuando lo llamo no sube.', '{}', 'PENDING', '{}', '{"Poll": null}', current_timestamp - interval '6 day' - interval '2 hours' - interval '56 minutes'),
       ('9bsv0s4ilem002n1mt9g', '9bsv0s4ilem002n1mt1g', '9bsv0s65sclg02m57a50', 'Es verdad, no va muy bien', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '6 days' - interval '1 hour' - interval '13 minutes'),
       ('9bsv0s4ilem002n1mta0', '9bsv0s4ilem002n1mt20', '9bsv0s65sclg02m57a4g', 'Anoche se escapó nuestro perro. Es un labrador blanco, todavía tiene 8 meses. Si lo véis por el barrio por favor avisadme a mí o a mi mujer. Muchísimas gracias.', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '2 days' - interval '7 hour' - interval '41 minutes'),
       ('9bsv0s4ilem002n1mtag', '9bsv0s4ilem002n1mt2g', '9bsv0s65sclg02m57a6g', 'No te preocupes Daniel, seguro que aparece.', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '6 days' + interval '1 hour' - interval '13 minutes'),
       ('9bsv0s4ilem002n1mtb0', '9bsv0s4ilem002n1mt30', '9bsv0s65sclg02m57a4g', 'Muchas gracias', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '6 days' + interval '1 hour' - interval '26 minutes'),
       ('9bsv0s4ilem002n1mtbg', '9bsv0s4ilem002n1mt3g', '9bsv0s65sclg02m57a9g', 'Nosotros no hemos visto nada. Si le vemos te avisamos', '{}', null, '{}', '{"Poll": null}', current_timestamp - interval '6 days' + interval '5 hour' - interval '49 minutes')
;
`

type SeedCommand struct {
	class.Command
	database database.Database
}

// TODO: CHANGE THIS, NEW SEED COMMAND SHOULD RETURN A CLI.COMMAND!
func NewSeedCommand(configuration internal.Configuration, logger core.Logger, database database.Database) *SeedCommand {
	return &SeedCommand{
		Command:  *class.NewCommand(configuration, logger),
		database: database,
	}
}

func (self *SeedCommand) Execute() error {
	err := self.database.Transaction(context.TODO(), func(ctx context.Context) error {
		numCommunities, err := self.database.Exec(ctx, communities)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		if numCommunities != 1 {
			return ErrGeneric().With("Community seed failed")
		}

		numUsers, err := self.database.Exec(ctx, users)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		if numUsers != 20 {
			return ErrGeneric().With("User seed failed")
		}

		numMemberships, err := self.database.Exec(ctx, memberships)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		if numMemberships != 20 {
			return ErrGeneric().With("Membership seed failed")
		}

		numPosts, err := self.database.Exec(ctx, posts)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		if numPosts != 16 {
			return ErrGeneric().With("Post seed failed")
		}

		numPostHistories, err := self.database.Exec(ctx, post_histories)
		if err != nil {
			return ErrGeneric().Wrap(err)
		}

		if numPostHistories != 16 {
			return ErrGeneric().With("Post History seed failed")
		}

		return nil
	})
	if err != nil {
		return ErrGeneric().Wrap(err)
	}

	return nil
}

// func (self *SeedCommand) Execute() *cli.Command {
// 	arguments := &payload.SeedArguments{}
// 	return self.Handle(class.CommandEndpoint{
// 		Name:        "seed",
// 		Description: "seed database with communities, users, memberships, posts and post histories",
// 		Arguments:   arguments,
// 	}, func(ctx *cli.Context) error {
// 		err := self.database.Transaction(context.TODO(), func(ctx context.Context) error {
// 			numCommunities, err := self.database.Exec(ctx, communities)
// 			if err != nil {
// 				return ErrGeneric().Wrap(err)
// 			}

// 			if numCommunities != 1 {
// 				return ErrGeneric().With("Community seed failed")
// 			}

// 			numUsers, err := self.database.Exec(ctx, users)
// 			if err != nil {
// 				return ErrGeneric().Wrap(err)
// 			}

// 			if numUsers != 20 {
// 				return ErrGeneric().With("User seed failed")
// 			}

// 			numMemberships, err := self.database.Exec(ctx, memberships)
// 			if err != nil {
// 				return ErrGeneric().Wrap(err)
// 			}

// 			if numMemberships != 20 {
// 				return ErrGeneric().With("Membership seed failed")
// 			}

// 			numPosts, err := self.database.Exec(ctx, posts)
// 			if err != nil {
// 				return ErrGeneric().Wrap(err)
// 			}

// 			if numPosts != 16 {
// 				return ErrGeneric().With("Post seed failed")
// 			}

// 			numPostHistories, err := self.database.Exec(ctx, post_histories)
// 			if err != nil {
// 				return ErrGeneric().Wrap(err)
// 			}

// 			if numPostHistories != 16 {
// 				return ErrGeneric().With("Post History seed failed")
// 			}

// 			return nil
// 		})
// 		if err != nil {
// 			return ErrGeneric().Wrap(err)
// 		}

// 		return nil
// 	})
// }
